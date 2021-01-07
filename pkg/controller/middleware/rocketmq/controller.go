package rocketmq

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appv1lister "k8s.io/client-go/listers/apps/v1"
	corev1lister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	crdclient "l0calh0st.cn/clickpaas-operator/pkg/client/clientset/versioned"
	"l0calh0st.cn/clickpaas-operator/pkg/client/clientset/versioned/scheme"
	middlewareinformer "l0calh0st.cn/clickpaas-operator/pkg/client/informers/externalversions"
	middlewarelister "l0calh0st.cn/clickpaas-operator/pkg/client/listers/middleware/v1alpha1"
	"l0calh0st.cn/clickpaas-operator/pkg/controller"
	"l0calh0st.cn/clickpaas-operator/pkg/operator"
	"l0calh0st.cn/clickpaas-operator/pkg/operator/middleware/rocketmq"
	"time"
)

type rocketmqController struct {
	controller.BaseController

	kubeClient kubernetes.Interface
	crdClient crdclient.Interface

	rocketmqLister middlewarelister.RocketmqLister

	deploymentLister appv1lister.DeploymentLister
	statefulSetLister appv1lister.StatefulSetLister
	serviceLister corev1lister.ServiceLister

	queue workqueue.RateLimitingInterface
	recorder record.EventRecorder
	cacheSyncedList []cache.InformerSynced

	operator operator.IOperator
}


func NewRocketmqController(kubeClient kubernetes.Interface, crdClient crdclient.Interface,
	kubeInformerFactory informers.SharedInformerFactory, midInformerFactory middlewareinformer.SharedInformerFactory)*rocketmqController{
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.V(2).Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events(corev1.NamespaceAll)})

	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "rocketmq"})

	return newRocketmqController(kubeClient, crdClient, kubeInformerFactory, midInformerFactory, recorder)
}


func newRocketmqController(kubeClient kubernetes.Interface, crdClient crdclient.Interface,
	kubeInformerFactory informers.SharedInformerFactory, midInformerFactory middlewareinformer.SharedInformerFactory,recorder record.EventRecorder,
	)*rocketmqController{

	controller := &rocketmqController{
		kubeClient: kubeClient,
		crdClient: crdClient,
		recorder: recorder,
		queue: workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}

	rocketmqInformer := midInformerFactory.Middleware().V1alpha1().Rocketmqs()
	rocketmqInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.onAdd,
		UpdateFunc: controller.onUpdate,
		DeleteFunc: controller.onDelete,
	})
	controller.rocketmqLister = rocketmqInformer.Lister()
	controller.cacheSyncedList = append(controller.cacheSyncedList, rocketmqInformer.Informer().HasSynced)

	deploymentInformer := kubeInformerFactory.Apps().V1().Deployments()
	controller.deploymentLister = deploymentInformer.Lister()
	controller.cacheSyncedList = append(controller.cacheSyncedList, deploymentInformer.Informer().HasSynced)

	statefulSetInformer := kubeInformerFactory.Apps().V1().StatefulSets()
	controller.statefulSetLister = statefulSetInformer.Lister()
	controller.cacheSyncedList = append(controller.cacheSyncedList, statefulSetInformer.Informer().HasSynced)

	controller.operator = rocketmq.NewRocketmqOperator(kubeClient, controller.rocketmqLister, controller.deploymentLister,
		controller.statefulSetLister, controller.serviceLister)

	return controller
}

func(c *rocketmqController)Start(ctx context.Context, threads int)error{
	if ok := cache.WaitForCacheSync(ctx.Done(), c.cacheSyncedList...); !ok {
		return fmt.Errorf("wait all informer cache has synced failed")
	}
	logrus.Infof("Rocketmq Controller has synced, ready to reconciling....")
	for i := 0; i< threads; i++{
		go wait.Until(c.runWorker, 5 * time.Second, ctx.Done())
	}
	<- ctx.Done()
	return ctx.Err()
}
func(c *rocketmqController)Stop(stopCh <- chan struct{})error{
	c.queue.ShutDown()
	<- stopCh
	return nil
}

func(c *rocketmqController)runWorker(){
	for c.processNextItem() {}
}

func(c *rocketmqController)processNextItem()bool{
	obj,shutdown := c.queue.Get()
	if shutdown{
		return false
	}
	err := func(obj interface{}) error {
		defer c.queue.Done(obj)
		key,ok := obj.(string)
		if !ok {
			runtime.HandleError(fmt.Errorf("except got string from workqueue, but got %#v", obj))
			return nil
		}
		if err := c.operator.Reconcile(key); err != nil{
			c.queue.AddRateLimited(key)
			return fmt.Errorf("reconcile failed requeue, %s:%s", key, err)
		}
		c.queue.Forget(key)
		logrus.Infof("reconcile successfully %s", key)
		return nil
	} (obj)
	if err != nil{
		runtime.HandleError(err)
	}
	return true
}

func(c *rocketmqController)onAdd(obj interface{}){
	rocket := obj.(*crdv1alpha1.Rocketmq)
	crdv1alpha1.WithDefaultsRocketmq(rocket)
	logrus.Infof("%d %d %d", rocket.Spec.FastPort, rocket.Spec.ListenPort, rocket.Spec.HaPort)
	c.recorder.Event(rocket, corev1.EventTypeNormal, RocketmqEventReasonOnAdded, eventMessage(rocket, RocketmqEventReasonOnAdded))
	for _, hook := range c.GetHooks(){
		hook.OnAdd(rocket)
	}
	c.enqueue(rocket)
}

func(c *rocketmqController)onDelete(obj interface{}){
	var rocketmq *crdv1alpha1.Rocketmq
	switch obj.(type) {
	case *crdv1alpha1.Rocketmq:
		rocketmq = obj.(*crdv1alpha1.Rocketmq)
	case cache.DeletedFinalStateUnknown:
		deleteObj := obj.(cache.DeletedFinalStateUnknown).Obj
		rocketmq = deleteObj.(*crdv1alpha1.Rocketmq)
	}
	if rocketmq == nil{
		return
	}

	c.recorder.Event(rocketmq, corev1.EventTypeNormal, RocketmqEventReasonOnDelete,eventMessage(rocketmq, RocketmqEventReasonOnDelete))
	for _, hook := range c.GetHooks(){
		hook.OnDelete(rocketmq)
	}
	// todo if do something else when delete, eg gc some resource or back some resource if necessary
}

func(c *rocketmqController)onUpdate(oldObj,newObj interface{}){
	oldRocketmq := oldObj.(*crdv1alpha1.Rocketmq)
	newRocketmq := newObj.(*crdv1alpha1.Rocketmq)
	if oldRocketmq.ResourceVersion == newRocketmq.ResourceVersion{
		// if the resourceVersion of oldObject and newObject are same, then do nothing and return
		return
	}
	for _,hook := range c.GetHooks(){
		hook.OnUpdate(newRocketmq)
	}
	c.enqueue(newRocketmq)
}

func(c *rocketmqController)enqueue(obj interface{}){
	key,err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil{
		runtime.HandleError(fmt.Errorf("cannot get key from queue, %s", err))
		return
	}
	c.queue.AddRateLimited(key)
}