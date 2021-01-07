package gcache

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
	"l0calh0st.cn/clickpaas-operator/pkg/operator/middleware/gcache"
	"context"
	"time"
)

type redisGCacheController struct {
	controller.BaseController
	kubeClient kubernetes.Interface
	crdClient crdclient.Interface

	gCacheLister crdlister.RedisGCacheLister
	statefulSetLister appv1lister.StatefulSetLister
	serviceLister corev1lister.ServiceLister

	queue workqueue.RateLimitingInterface
	cacheSyncedList []cache.InformerSynced
	recorder record.EventRecorder

	operator operator.IOperator
}

func NewRedisGCacheController(kubeClient kubernetes.Interface, crdClient crdclient.Interface,
	kubeInformerFactory informers.SharedInformerFactory, crdInformerFactory crdinformer.SharedInformerFactory)*redisGCacheController{

	eventBroadCaster := record.NewBroadcaster()
	eventBroadCaster.StartLogging(glog.V(2).Infof)
	eventBroadCaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events(corev1.NamespaceAll)})
	recorder := eventBroadCaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "redis-gcache"})

	return newRedisGCacheController(kubeClient, crdClient, kubeInformerFactory, crdInformerFactory, recorder)

}


func newRedisGCacheController(kubeClient kubernetes.Interface, crdClient crdclient.Interface,
	kubeInformerFactory informers.SharedInformerFactory, crdInformerFactory crdinformer.SharedInformerFactory, recorder record.EventRecorder)*redisGCacheController{

	controller := &redisGCacheController{
		kubeClient:        kubeClient,
		crdClient: crdClient,
		queue:             workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		cacheSyncedList:   []cache.InformerSynced{},
		recorder:          recorder,
	}

	redisGCacheInformer := crdInformerFactory.Middleware().V1alpha1().RedisGCaches()
	redisGCacheInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.onAdd,
		UpdateFunc: controller.onUpdate,
		DeleteFunc: controller.onDelete,
	})
	controller.gCacheLister = redisGCacheInformer.Lister()
	controller.cacheSyncedList = append(controller.cacheSyncedList, redisGCacheInformer.Informer().HasSynced)

	statefulSetInformer := kubeInformerFactory.Apps().V1().StatefulSets()
	controller.statefulSetLister = statefulSetInformer.Lister()
	controller.cacheSyncedList = append(controller.cacheSyncedList, statefulSetInformer.Informer().HasSynced)

	serviceInformer := kubeInformerFactory.Core().V1().Services()
	controller.serviceLister = serviceInformer.Lister()
	controller.cacheSyncedList = append(controller.cacheSyncedList, serviceInformer.Informer().HasSynced)

	controller.operator = gcache.NewRedisGCacheOperator(kubeClient, controller.gCacheLister, controller.statefulSetLister, controller.serviceLister)
	return controller
}

func(c *redisGCacheController)Start(ctx context.Context, threads int)error{
	if ok := cache.WaitForCacheSync(ctx.Done(), c.cacheSyncedList...); !ok {
		return fmt.Errorf("wait all informer has synced failed")
	}
	for i := 0 ; i< threads ;i++{
		go wait.Until(c.runWorker, 5 * time.Second, ctx.Done())
	}
	<- ctx.Done()
	return ctx.Err()
}

func(c *redisGCacheController)Stop(stopCh <- chan struct{})error{
	c.queue.ShutDown()
	return nil
}


func(c *redisGCacheController)runWorker(){
	defer runtime.HandleCrash()
	for c.processNextItem(){}
}

func(c *redisGCacheController)processNextItem()bool{
	key, shutdown := c.queue.Get()
	if shutdown{
		runtime.HandleError(fmt.Errorf("get key from workqueue failed"))
		return false
	}
	err := func(obj interface{})error {
		defer c.queue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.queue.Forget(obj)
			runtime.HandleError(fmt.Errorf("except got string from workqueue, but got %#v", obj))
			return nil
		}
		if err := c.operator.Reconcile(key); err != nil{
			c.queue.AddRateLimited(key)
			return fmt.Errorf("reconcile redisGCache failed: %v", err)
		}
		c.queue.Forget(obj)
		logrus.Infof("Reconcole redisGCache successfully")
		return nil
	}(key)
	if err != nil{
		runtime.HandleError(err)
	}
	return true

}

func(c *redisGCacheController)onAdd(obj interface{}){
	redis := obj.(*crdv1alpha1.RedisGCache)
	crdv1alpha1.WithDefaultsRedisGCache(redis)
	c.recorder.Event(redis, corev1.EventTypeNormal, RedisGCacheEventReasonOnAdded, fmt.Sprintf("'%v:%v' added", redis.GetName(), redis.GetNamespace()))
	for _, hook := range c.GetHooks(){
		hook.OnAdd(redis)
	}
	c.enqueue(redis)
}

func(c *redisGCacheController)onUpdate(oldObj,newObj interface{}){
	oldRedis := oldObj.(*crdv1alpha1.RedisGCache)
	newRedis := newObj.(*crdv1alpha1.RedisGCache)
	if oldRedis.ResourceVersion == newRedis.ResourceVersion{
		return
	}
	c.recorder.Event(newRedis, corev1.EventTypeNormal, RedisGCacheEventReasonOnUpdate, fmt.Sprintf("'%v:%v' updated", newRedis.GetName(), newRedis.GetNamespace()))
	for _, hook := range c.GetHooks(){
		hook.OnUpdate(newRedis)
	}
	c.enqueue(newRedis)
}

func(c *redisGCacheController)onDelete(obj interface{}){
	var redis *crdv1alpha1.RedisGCache
	switch obj.(type) {
	case *crdv1alpha1.RedisGCache:
		redis = obj.(*crdv1alpha1.RedisGCache)
	case cache.DeletedFinalStateUnknown:
		deleteObj := obj.(cache.DeletedFinalStateUnknown).Obj
		redis = deleteObj.(*crdv1alpha1.RedisGCache)
	}
	if redis == nil{
		return
	}
	c.recorder.Event(redis, corev1.EventTypeNormal, RedisGCacheEventReasonOnDelete, fmt.Sprintf("'%v:%v' delete",redis.GetName(), redis.GetNamespace()))
	for _, hook := range c.GetHooks(){
		hook.OnDelete(redis)
	}
}

func(c *redisGCacheController)enqueue(obj interface{}){
	key,err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil{
		logrus.Errorf("faile to get key for %v", obj)
		return
	}
	c.queue.AddRateLimited(key)
}