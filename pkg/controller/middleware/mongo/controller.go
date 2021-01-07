package mongo

import (
	"context"
	"fmt"
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
	"k8s.io/klog/v2"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	crdclient "l0calh0st.cn/clickpaas-operator/pkg/client/clientset/versioned"
	"l0calh0st.cn/clickpaas-operator/pkg/client/clientset/versioned/scheme"
	crdinformer "l0calh0st.cn/clickpaas-operator/pkg/client/informers/externalversions"
	"l0calh0st.cn/clickpaas-operator/pkg/client/listers/middleware/v1alpha1"
	"l0calh0st.cn/clickpaas-operator/pkg/controller"
	"l0calh0st.cn/clickpaas-operator/pkg/operator"
	"l0calh0st.cn/clickpaas-operator/pkg/operator/middleware/mongo"
	"time"
)

type mongoController struct {
	controller.BaseController
	kubeClient kubernetes.Interface
	crdClient crdclient.Interface

	mongoLister v1alpha1.MongoClusterLister
	statefulSetLister appv1lister.StatefulSetLister
	serviceLister corev1lister.ServiceLister

	queue workqueue.RateLimitingInterface
	recorder record.EventRecorder
	cacheSyncedList []cache.InformerSynced

	operator operator.IOperator
}

func NewMongoController(kubeClient kubernetes.Interface, crdClient crdclient.Interface,
	kubeInformerFactory informers.SharedInformerFactory, crdInformerFactory crdinformer.SharedInformerFactory)*mongoController{
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.V(2).Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events(corev1.NamespaceAll)})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "mongo-controller"})
	return newMongoController(kubeClient, crdClient, kubeInformerFactory, crdInformerFactory, recorder)
}

func newMongoController(kubeClient kubernetes.Interface, crdClient crdclient.Interface, kubeInformerFactory informers.SharedInformerFactory,
	crdInformerFactory crdinformer.SharedInformerFactory, recorder record.EventRecorder)*mongoController{
	controller := &mongoController{
		kubeClient:        kubeClient,
		crdClient:         crdClient,
		queue:             workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		recorder:          recorder,
		cacheSyncedList:   []cache.InformerSynced{},
	}
	mongoInformer := crdInformerFactory.Middleware().V1alpha1().MongoClusters()
	mongoInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.onAdd,
		UpdateFunc: controller.onUpdate,
		DeleteFunc: controller.onDelete,
	})
	controller.mongoLister = mongoInformer.Lister()
	controller.cacheSyncedList = append(controller.cacheSyncedList, mongoInformer.Informer().HasSynced)

	statefulSetInformer := kubeInformerFactory.Apps().V1().StatefulSets()
	controller.statefulSetLister = statefulSetInformer.Lister()
	controller.cacheSyncedList = append(controller.cacheSyncedList, statefulSetInformer.Informer().HasSynced)

	serviceInformer := kubeInformerFactory.Core().V1().Services()
	controller.serviceLister = serviceInformer.Lister()
	controller.cacheSyncedList = append(controller.cacheSyncedList, serviceInformer.Informer().HasSynced)

	controller.operator = mongo.NewMongoOperator(kubeClient, controller.mongoLister, controller.serviceLister, controller.statefulSetLister)

	return controller
}


func(c *mongoController)Start(ctx context.Context, threads int)error{
	// wait all informer cache synced
	if ok := cache.WaitForCacheSync(ctx.Done(), c.cacheSyncedList ...); !ok {
		return fmt.Errorf("wait all informer cache has synced failed")
	}
	logrus.Infof("Mongo Controller has started, ready to reconciling")
	for i:= 0 ;i < threads; i++{
		go wait.Until(c.runWorker, 5 * time.Second, ctx.Done())
	}
	<- ctx.Done()
	return ctx.Err()
}

func(c *mongoController)runWorker(){
	runtime.HandleCrash()
	for c.processNextItem(){}
}

func (c *mongoController)processNextItem()bool{
	obj,shutdown := c.queue.Get()
	if shutdown{
		return false
	}
	err := func(obj interface{}) error{
		defer c.queue.Done(obj)
		var key string
		var ok bool
		if key,ok = obj.(string); !ok {
			runtime.HandleError(fmt.Errorf("except got string from workqueue, but got %#v", key))
			return nil
		}
		if err := c.operator.Reconcile(key); err != nil{
			c.queue.AddRateLimited(key)
			return fmt.Errorf("reconcile mongo cluster failed: '%v':%v",key, err)
		}
		c.queue.Forget(key)
		logrus.Infof("reconcile mongo cluster successfully")
		return nil
	}(obj)
	if err != nil{
		runtime.HandleError(err)
		return true
	}
	return true

}

func(c *mongoController)Stop(stopCh <- chan struct{})error{
	c.queue.ShutDown()
	return nil
}

func(c *mongoController)onAdd(obj interface{}){
	mongo := obj.(*crdv1alpha1.MongoCluster)
	crdv1alpha1.WithDefaultsMongoCluster(mongo)
	for _, hook := range c.GetHooks(){
		hook.OnAdd(mongo)
	}
	c.enqueue(mongo)
}

func(c *mongoController)onUpdate(oldObj,newObj interface{}){
	oldMongo := oldObj.(*crdv1alpha1.MongoCluster)
	newMongo := newObj.(*crdv1alpha1.MongoCluster)
	if oldMongo.ResourceVersion == newMongo.ResourceVersion{
		// resource version are same, then do nothing, and return
		return
	}
	for _,hook := range c.GetHooks(){
		hook.OnUpdate(newMongo)
	}
	c.enqueue(newMongo)

}

func(c *mongoController)onDelete(obj interface{}){
	var mongo *crdv1alpha1.MongoCluster
	switch obj.(type) {
	case *crdv1alpha1.MongoCluster:
		mongo = obj.(*crdv1alpha1.MongoCluster)
	case cache.DeletedFinalStateUnknown:
		deleteObj := obj.(cache.DeletedFinalStateUnknown).Obj
		mongo = deleteObj.(*crdv1alpha1.MongoCluster)
	}
	if mongo == nil{
		return
	}
	for _, hook := range c.GetHooks(){
		hook.OnDelete(mongo)
	}
}

func(c *mongoController)enqueue(obj interface{}){
	key,err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil{
		logrus.Errorf("failed to get key from workqueue %v", err)
	}
	c.queue.AddRateLimited(key)
}