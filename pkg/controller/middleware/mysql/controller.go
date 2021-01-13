package mysql

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
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
	"l0calh0st.cn/clickpaas-operator/pkg/operator/middleware/mysql"
	"time"

	appv1lister "k8s.io/client-go/listers/apps/v1"
	corev1lister "k8s.io/client-go/listers/core/v1"
)

type mysqlController struct {
	controller.BaseController
	kubeClient kubernetes.Interface
	crdClient crdclient.Interface

	cacheSyncedList []cache.InformerSynced
	queue workqueue.RateLimitingInterface

	statefulSetLister appv1lister.StatefulSetLister
	serviceLister corev1lister.ServiceLister
	mysqlClusterLister crdlister.MysqlClusterLister

	recorder record.EventRecorder

	operator operator.IOperator

}

func NewMysqlController(
	kubeClient kubernetes.Interface,
	crdClient crdclient.Interface,
	crdInformerFactory crdinformer.SharedInformerFactory,
	kubeInformerFactory informers.SharedInformerFactory)*mysqlController{

	eventBroadCaster := record.NewBroadcaster()
	eventBroadCaster.StartLogging(glog.V(2).Infof)
	eventBroadCaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events(metav1.NamespaceAll)})
	recorder := eventBroadCaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "mysql-controller"})

	return newMysqlClusterController(kubeClient, crdClient, recorder, kubeInformerFactory, crdInformerFactory)
}


func newMysqlClusterController(
	kubeClient kubernetes.Interface,
	crdClient crdclient.Interface,
	recorder record.EventRecorder,
	kubeInformerFactory informers.SharedInformerFactory,
	crdInformerFactory crdinformer.SharedInformerFactory,
	)*mysqlController{
	controller := &mysqlController{
		kubeClient:         kubeClient,
		crdClient:          crdClient,
		recorder: recorder,
		queue: workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		cacheSyncedList: []cache.InformerSynced{},
	}

	mysqlInformer := crdInformerFactory.Middleware().V1alpha1().MysqlClusters()
	mysqlInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.onAdd,
		UpdateFunc: controller.onUpdate,
		DeleteFunc: controller.onDelete,
	})
	controller.cacheSyncedList = append(controller.cacheSyncedList, mysqlInformer.Informer().HasSynced)
	controller.mysqlClusterLister = mysqlInformer.Lister()

	svcInformer := kubeInformerFactory.Core().V1().Services()
	controller.cacheSyncedList = append(controller.cacheSyncedList, svcInformer.Informer().HasSynced)
	controller.serviceLister = svcInformer.Lister()

	statefulSetInformer := kubeInformerFactory.Apps().V1().StatefulSets()
	controller.cacheSyncedList = append(controller.cacheSyncedList, statefulSetInformer.Informer().HasSynced)
	controller.statefulSetLister = statefulSetInformer.Lister()

	controller.operator = mysql.NewMysqlClusterOperator(kubeClient,controller.mysqlClusterLister,controller.statefulSetLister, controller.serviceLister)

	return controller
}


// Start run controller
func (c *mysqlController) Start(ctx context.Context, controllerThreads int)error{
	logrus.Warn("begin start controller")
	if ok := cache.WaitForCacheSync(ctx.Done(), c.cacheSyncedList...); !ok{
		return fmt.Errorf("timeout wait for cache synced")
	}
	logrus.Infof("Mysql Controller has started, ready to reclioning")
	for i := 0 ;i < controllerThreads; i++{
		go wait.Until(c.runWorker, time.Second, ctx.Done())
	}
	<- ctx.Done()
	return ctx.Err()
}

// Stop shutdown workQueue
func (c *mysqlController)Stop(stopCh <- chan struct{})error{
	logrus.Info("Stopping the hdfs controller")
	c.queue.ShutDown()
	<- stopCh
	return nil
}
func(c *mysqlController)runWorker(){
	defer runtime.HandleCrash()
	for c.processNextItem() {}
}


func(c *mysqlController) processNextItem()bool{
	obj,shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	err := func(obj interface{}) error{
		defer c.queue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok{
			c.queue.Forget(obj)
			runtime.HandleError(fmt.Errorf("except got string in workqueue, but got %#v", obj))
			return nil
		}
		if err := c.operator.Reconcile(key); err != nil{
			c.queue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%v':%v, requeue",key ,err)
		}
		c.queue.Forget(obj)
		logrus.Infof("successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil{
		runtime.HandleError(err)
		return true
	}
	return true
}

func(c *mysqlController)enqueue(obj interface{}){
	key,err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil{
		logrus.Errorf("failed to get key for %v:%v", obj, err)
		return
	}
	c.queue.AddRateLimited(key)
}



func(c *mysqlController)onAdd(obj interface{}){
	mysql := obj.(*crdv1alpha1.MysqlCluster)
	crdv1alpha1.WithDefaultsMysqlCluster(mysql)
	c.recorder.Eventf(mysql, corev1.EventTypeNormal, string("Created"), "%v created", mysql.GetName())
	for _,hook := range c.GetHooks(){
		hook.OnAdd(obj)
	}
	c.enqueue(mysql)
}

func(c *mysqlController)onUpdate(oldObj,newObj interface{}){
	oldCluster := oldObj.(*crdv1alpha1.MysqlCluster)
	newCluster := newObj.(*crdv1alpha1.MysqlCluster)
	if oldCluster.ResourceVersion == newCluster.ResourceVersion {
		return
	}
	//if !equality.Semantic.DeepEqual(oldCluster.Spec, newCluster.Spec){
	//	// if semantic change ,then delete old one, and crate newOne
	//	err := c.crdClient.MiddlewareV1alpha1().MysqlClusters(oldCluster.GetNamespace()).Delete(context.TODO(), oldCluster.GetName(), metav1.DeleteOptions{})
	//	if err != nil{
	//
	//	}
	//	newCluster,err = c.crdClient.MiddlewareV1alpha1().MysqlClusters(newCluster.GetNamespace()).Create(context.TODO(), newCluster, metav1.CreateOptions{})
	//}
	for _,hook := range c.GetHooks(){
		hook.OnUpdate(newObj)
	}
	c.enqueue(newCluster)
}


func(c *mysqlController)onDelete(obj interface{}){
	// todo do some thing when delete, eg backup or other some thing if necessary
	var mysql *crdv1alpha1.MysqlCluster
	switch obj.(type) {
	case *crdv1alpha1.MysqlCluster:
		mysql = obj.(*crdv1alpha1.MysqlCluster)
	case cache.DeletedFinalStateUnknown:
		deleteObj := obj.(cache.DeletedFinalStateUnknown).Obj
		mysql = deleteObj.(*crdv1alpha1.MysqlCluster)
	}
	if mysql != nil{
	}
	for _,hook := range c.GetHooks(){
		hook.OnDelete(obj)
	}
}
