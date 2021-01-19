package zookeeper

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appv1lister "k8s.io/client-go/listers/apps/v1"
	corev1lister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	crdclient "l0calh0st.cn/clickpaas-operator/pkg/client/clientset/versioned"
	"l0calh0st.cn/clickpaas-operator/pkg/client/clientset/versioned/scheme"
	middlewareinformer "l0calh0st.cn/clickpaas-operator/pkg/client/informers/externalversions"
	middlewarelister "l0calh0st.cn/clickpaas-operator/pkg/client/listers/middleware/v1alpha1"
	"l0calh0st.cn/clickpaas-operator/pkg/controller"
	"l0calh0st.cn/clickpaas-operator/pkg/controller/middleware/zookeeper/handler"
	"l0calh0st.cn/clickpaas-operator/pkg/operator"
	corev1 "k8s.io/api/core/v1"
	"l0calh0st.cn/clickpaas-operator/pkg/operator/middleware/zookeeper"
	"time"
)

type zookeeperController struct {
	controller.BaseController
	kubeClient kubernetes.Interface
	crdClient crdclient.Interface

	zkLister middlewarelister.ZookeeperClusterLister
	statefulSetLister appv1lister.StatefulSetLister
	serviceLister corev1lister.ServiceLister
	configMapLister corev1lister.ConfigMapLister
	podLister corev1lister.PodLister


	queue workqueue.RateLimitingInterface
	recorder record.EventRecorder
	cacheSyncedList []cache.InformerSynced

	operator operator.IOperator
}


func NewZookeeperController(kubeClient kubernetes.Interface, crdClient crdclient.Interface,restConfig *rest.Config,
	middlewareInformerFactory middlewareinformer.SharedInformerFactory, kubeInformerFactory informers.SharedInformerFactory,
	)*zookeeperController{

	eventBroadCaster := record.NewBroadcaster()
	eventBroadCaster.StartLogging(glog.V(2).Infof)
	eventBroadCaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events(corev1.NamespaceAll)})

	recorder := eventBroadCaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "zookeeper-controller"})
	return newZookeeperController(kubeClient, crdClient,restConfig, middlewareInformerFactory, kubeInformerFactory, recorder)
}


func newZookeeperController(kubeClient kubernetes.Interface, crdClient crdclient.Interface,restConfig *rest.Config,
	middleInformerFactory middlewareinformer.SharedInformerFactory, kubeInformerFactory informers.SharedInformerFactory,
	recorder record.EventRecorder)*zookeeperController{

	controller := &zookeeperController{
		kubeClient: kubeClient,
		crdClient: crdClient,
		queue: workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		cacheSyncedList: []cache.InformerSynced{},
		recorder: recorder,
	}
	zkInformer := middleInformerFactory.Middleware().V1alpha1().ZookeeperClusters()
	zkInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.onAdd,
		UpdateFunc: controller.onUpdate,
		DeleteFunc: controller.onDelete,
	})
	controller.zkLister = zkInformer.Lister()
	controller.cacheSyncedList = append(controller.cacheSyncedList, zkInformer.Informer().HasSynced)

	serviceInformer := kubeInformerFactory.Core().V1().Services()
	serviceInformer.Informer().AddEventHandler(handler.NewServiceEventHandler(controller.zkLister, controller.enqueue))
	controller.serviceLister = serviceInformer.Lister()
	controller.cacheSyncedList = append(controller.cacheSyncedList, serviceInformer.Informer().HasSynced)

	statefulSetInformer := kubeInformerFactory.Apps().V1().StatefulSets()
	controller.statefulSetLister = statefulSetInformer.Lister()
	controller.cacheSyncedList = append(controller.cacheSyncedList, statefulSetInformer.Informer().HasSynced)

	podInformer := kubeInformerFactory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(handler.NewPodEventHandler(controller.zkLister, controller.enqueue))
	controller.cacheSyncedList = append(controller.cacheSyncedList, podInformer.Informer().HasSynced)
	controller.podLister = podInformer.Lister()

	configMapInformer := kubeInformerFactory.Core().V1().ConfigMaps()
	configMapInformer.Informer().AddEventHandler(handler.NewConfigMapEventHandler(controller.zkLister, controller.enqueue))
	controller.cacheSyncedList = append(controller.cacheSyncedList, configMapInformer.Informer().HasSynced)
	controller.configMapLister = configMapInformer.Lister()

	controller.operator = zookeeper.NewZookeeperOperator(kubeClient,restConfig ,controller.zkLister, controller.statefulSetLister,
		controller.serviceLister, controller.configMapLister, controller.podLister)
	return controller
}

func(c *zookeeperController)Start(ctx context.Context, threads int)error{
	if ok := cache.WaitForCacheSync(ctx.Done(), c.cacheSyncedList...); !ok{
		return fmt.Errorf("wait for all informer cache synced failed")
	}
	logrus.Infof("zookeeper controller has synced, ready to reconciling.....")
	for i := 0 ; i< threads; i++{
		go wait.Until(c.runWorker, 5 * time.Second, ctx.Done())
	}
	<- ctx.Done()
	return ctx.Err()
}
func(c *zookeeperController)Stop(stopCh <- chan struct{})error{
	c.queue.ShutDown()
	<- stopCh
	return nil
}

func(c *zookeeperController)runWorker(){
	for c.processNextItem(){}
}
func(c *zookeeperController)processNextItem()bool{
	item,shutdown := c.queue.Get()
	if shutdown{
		return false
	}
	err := func()error {
		defer c.queue.Done(item)
		key,ok := item.(string)
		if !ok {
			runtime.HandleError(fmt.Errorf("except got string from workqueue, but got %#v", item))
			return nil
		}
		if err := c.operator.Reconcile(key); err != nil{
			c.queue.AddRateLimited(key)
			return fmt.Errorf("reconciling %s failed, %s", key, err)
		}
		logrus.Infof("successfully synced '%v'", key)
		c.queue.Forget(key)
		return nil
	}()
	if err != nil{
		runtime.HandleError(err)
	}
	return true
}

func(c *zookeeperController)onAdd(obj interface{}){
	zk := obj.(*crdv1alpha1.ZookeeperCluster)
	crdv1alpha1.WithDefaultsZookeeper(zk)
	c.recorder.Event(zk, corev1.EventTypeNormal, ZookeeperEventReasonOnAdded, eventMessage(zk, ZookeeperEventReasonOnAdded))
	for _, hook := range c.GetHooks(){
		hook.OnAdd(zk)
	}
	c.enqueue(zk)
}

func(c *zookeeperController)onDelete(obj interface{}){
	var zk *crdv1alpha1.ZookeeperCluster
	switch obj.(type) {
	case *crdv1alpha1.ZookeeperCluster:
		zk = obj.(*crdv1alpha1.ZookeeperCluster)
	case cache.DeletedFinalStateUnknown:
		deleteObj := obj.(cache.DeletedFinalStateUnknown).Obj
		zk = deleteObj.(*crdv1alpha1.ZookeeperCluster)
	}
	if zk == nil{
		return
	}
	c.recorder.Event(zk, corev1.EventTypeNormal, ZookeeperEventReasonOnDelete, eventMessage(zk,ZookeeperEventReasonOnDelete))
	for _,hook := range c.GetHooks(){
		hook.OnDelete(zk)
	}
	// todo if necessary
}

func(c *zookeeperController)onUpdate(oldObj,newObj interface{}){
	oldZk := oldObj.(*crdv1alpha1.ZookeeperCluster)
	newZk := newObj.(*crdv1alpha1.ZookeeperCluster)

	if oldZk.ResourceVersion == newZk.ResourceVersion{
		return
	}
	c.recorder.Event(newZk, corev1.EventTypeNormal, ZookeeperEventReasonOnUpdate, eventMessage(newZk, ZookeeperEventReasonOnUpdate))
	for _,hook := range c.GetHooks(){
		hook.OnUpdate(newZk)
	}
	c.enqueue(newZk)
}

func(c *zookeeperController)enqueue(obj interface{}){
	if key,err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj);err != nil{
		runtime.HandleError(fmt.Errorf("cannot get key from queue %s", err))
	}else {
		c.queue.Add(key)
	}
}


