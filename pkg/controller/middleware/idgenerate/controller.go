package idgenerate

import (
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
	crdinformer "l0calh0st.cn/clickpaas-operator/pkg/client/informers/externalversions"
	crdlister "l0calh0st.cn/clickpaas-operator/pkg/client/listers/middleware/v1alpha1"
	"l0calh0st.cn/clickpaas-operator/pkg/controller"
	"l0calh0st.cn/clickpaas-operator/pkg/operator"
	"context"
	"l0calh0st.cn/clickpaas-operator/pkg/operator/middleware/idgenerator"
	"time"
)

type idGeneratorController struct {
	controller.BaseController
	kubeClient kubernetes.Interface
	crdClient crdclient.Interface

	redisIdGeneratorLister crdlister.IdGenerateLister
	statefulSetLister appv1lister.StatefulSetLister
	serviceLister corev1lister.ServiceLister

	queue workqueue.RateLimitingInterface
	cacheSyncedList []cache.InformerSynced
	recorder record.EventRecorder

	operator operator.IOperator
}

func NewRedisIdGeneratorController(kubeClient kubernetes.Interface, crdClient crdclient.Interface,
	kubeInformerFactory informers.SharedInformerFactory, crdInformerFactory crdinformer.SharedInformerFactory,
)*idGeneratorController{
	eventBroadCaster := record.NewBroadcaster()
	eventBroadCaster.StartLogging(glog.V(2).Infof)
	eventBroadCaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events(corev1.NamespaceAll)})
	recorder := eventBroadCaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "idgenerator-controller"})

	return newRedisIdGeneratorController(kubeClient, crdClient, kubeInformerFactory, crdInformerFactory, recorder)
}

func newRedisIdGeneratorController(kubeClient kubernetes.Interface, crdClient crdclient.Interface,
	kubeInformerFactory informers.SharedInformerFactory, crdInformerFactory crdinformer.SharedInformerFactory,
	recorder record.EventRecorder)*idGeneratorController{

	controller := &idGeneratorController{
		kubeClient:             kubeClient,
		crdClient:              crdClient,
		queue:                  workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		cacheSyncedList:        []cache.InformerSynced{},
		recorder:               recorder,
	}

	idGeneratorInformer := crdInformerFactory.Middleware().V1alpha1().IdGenerates()
	idGeneratorInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.onAdd,
		UpdateFunc: controller.onUpdate,
		DeleteFunc: controller.onDelete,
	})
	controller.redisIdGeneratorLister = idGeneratorInformer.Lister()
	controller.cacheSyncedList = append(controller.cacheSyncedList, idGeneratorInformer.Informer().HasSynced)

	serviceInformer := kubeInformerFactory.Core().V1().Services()
	controller.serviceLister = serviceInformer.Lister()
	controller.cacheSyncedList = append(controller.cacheSyncedList, serviceInformer.Informer().HasSynced)

	statefulSetInformer := kubeInformerFactory.Apps().V1().StatefulSets()
	controller.cacheSyncedList = append(controller.cacheSyncedList, statefulSetInformer.Informer().HasSynced)

	controller.operator = idgenerator.NewIdGeneratorOperator(kubeClient, controller.redisIdGeneratorLister, controller.serviceLister, controller.statefulSetLister)
	return controller
}

func(c *idGeneratorController)Start(ctx context.Context, threads int)error{
	if ok := cache.WaitForCacheSync(ctx.Done(), c.cacheSyncedList...);!ok {
		return fmt.Errorf("waitForCacheSync failed")
	}
	logrus.Infof("IdGenerate Controller has started, ready to reconciling")
	for i := 0; i< threads; i++{
		go wait.Until(c.runWorker, 5 * time.Second, ctx.Done())
	}
	<- ctx.Done()
	return ctx.Err()
}
func(c *idGeneratorController)Stop(stopCh <- chan struct{})error{
	c.queue.ShutDown()
	<- stopCh
	return nil
}

func (c *idGeneratorController)runWorker(){
	for c.processNextItem() {}
}
func (c *idGeneratorController)processNextItem()bool{
	obj, shutdown := c.queue.Get()
	if shutdown{
		return false
	}
	err := func(obj interface{})error {
		defer c.queue.Done(obj)
		key,ok := obj.(string)
		if !ok {
			runtime.HandleError(fmt.Errorf("except got string from workqueue, but got %#v", key))
			return nil
		}
		if err := c.operator.Reconcile(key); err != nil{
			c.queue.AddRateLimited(key)
			return fmt.Errorf("reconcile failed: %s", err)
		}
		c.queue.Forget(obj)
		logrus.Infof("reconcile successfully")
		return nil
	}(obj)
	if err != nil{
		runtime.HandleError(err)
	}
	return true
}


func(c *idGeneratorController)onAdd(obj interface{}){
	idGenerator := obj.(*crdv1alpha1.IdGenerate)
	crdv1alpha1.WithDefaultsRedisIdGenerate(idGenerator)
	c.recorder.Event(idGenerator, corev1.EventTypeNormal, IdGeneratorEventReasonOnAdded, eventMessage(idGenerator, IdGeneratorEventReasonOnAdded))
	for _,hook := range c.GetHooks(){
		hook.OnAdd(idGenerator)
	}
	c.enqueue(idGenerator)

}
func(c *idGeneratorController)onUpdate(oldObj,newObj interface{}){
	oldGen := oldObj.(*crdv1alpha1.IdGenerate)
	newGen := newObj.(*crdv1alpha1.IdGenerate)
	if oldGen.ResourceVersion == newGen.ResourceVersion{
		logrus.Warnf("'%s:%s resourceVersion has not changed, skip recloncile", newGen.GetName(), newGen.GetNamespace())
		return
	}
	c.recorder.Event(newGen, corev1.EventTypeNormal, IdGeneratorEventReasonOnUpdate, eventMessage(newGen, IdGeneratorEventReasonOnUpdate))
	for _, hook := range c.GetHooks(){
		hook.OnUpdate(newGen)
	}
	c.enqueue(newGen)
}
func(c *idGeneratorController)onDelete(obj interface{}){
	var idgenerate *crdv1alpha1.IdGenerate
	switch obj.(type) {
	case *crdv1alpha1.IdGenerate:
		idgenerate = obj.(*crdv1alpha1.IdGenerate)
	case cache.DeletedFinalStateUnknown:
		deleteObj := obj.(cache.DeletedFinalStateUnknown).Obj
		idgenerate = deleteObj.(*crdv1alpha1.IdGenerate)
	}
	if idgenerate == nil{
		return
	}

	c.recorder.Event(idgenerate, corev1.EventTypeNormal, IdGeneratorEventReasonOnDelete, eventMessage(idgenerate, IdGeneratorEventReasonOnDelete))
	for _,hook := range c.GetHooks(){
		hook.OnDelete(idgenerate)
	}

}

func(c *idGeneratorController)enqueue(obj interface{}){
	key,err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(fmt.Errorf("failed to get key from workqueue: %s", err))
		return
	}
	c.queue.AddRateLimited(key)
}