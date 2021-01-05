package mysql

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
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
	"l0calh0st.cn/clickpaas-operator/pkg/operator/mysql"
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
	if err := crdv1alpha1.AddToScheme(scheme.Scheme); err != nil{
		logrus.Panic("NewMysqlController failed, register failed %v",err)
	}
	eventBroadCaster := record.NewBroadcaster()
	eventBroadCaster.StartLogging(glog.V(2).Infof)
	eventBroadCaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events(metav1.NamespaceAll)})
	recorder := eventBroadCaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "mysql"})

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
		cacheSyncedList: []cache.InformerSynced{},
	}

	mysqlInformer := crdInformerFactory.Middleware().V1alpha1().MysqlClusters()
	mysqlInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.onAdd,
		UpdateFunc: controller.onUpdate,
		DeleteFunc: controller.onDelete,
	})

	controller.cacheSyncedList = append(controller.cacheSyncedList, mysqlInformer.Informer().HasSynced)
	mysqlLister := mysqlInformer.Lister()

	svcInformer := kubeInformerFactory.Core().V1().Services()
	controller.cacheSyncedList = append(controller.cacheSyncedList, svcInformer.Informer().HasSynced)
	svcLister := svcInformer.Lister()

	statefulSetInformer := kubeInformerFactory.Apps().V1().StatefulSets()
	controller.cacheSyncedList = append(controller.cacheSyncedList, statefulSetInformer.Informer().HasSynced)
	ssLister := statefulSetInformer.Lister()

	controller.operator = mysql.NewMysqlClusterOperator(kubeClient,mysqlLister ,ssLister, svcLister)

	return controller
}


func(c *mysqlController)onAdd(obj interface{}){
	mysql := obj.(*crdv1alpha1.MysqlCluster)
	crdv1alpha1.WithDefaultsMysqlCluster(mysql)
	logrus.Info("Hdfs Cluster %s was added,enqueue for next handling")
	c.recorder.Eventf(mysql, corev1.EventTypeNormal, string("Created"), "%v", mysql.GetName())
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
	if !equality.Semantic.DeepEqual(oldCluster.Spec, newCluster.Spec){
		// force update hdfs cluster
	}
	for _,hook := range c.GetHooks(){
		hook.OnUpdate(newObj)
	}
	c.enqueue(newCluster)
}


func(c *mysqlController)onDelete(obj interface{}){
	var mysql *crdv1alpha1.MysqlCluster
	switch obj.(type) {
	case *crdv1alpha1.MysqlCluster:
		mysql = obj.(*crdv1alpha1.MysqlCluster)
	case cache.DeletedFinalStateUnknown:
		deleteObj := obj.(cache.DeletedFinalStateUnknown).Obj
		mysql = deleteObj.(*crdv1alpha1.MysqlCluster)
	}
	if mysql != nil{
		// 执行删除处理
	}
	for _,hook := range c.GetHooks(){
		hook.OnDelete(obj)
	}
}



func (c *mysqlController) Start(ctx context.Context, controllerThreads int)error{
	logrus.Warn("start controller")
	if ok := cache.WaitForCacheSync(ctx.Done(), c.cacheSyncedList...); !ok{
		return fmt.Errorf("timeout wait for cache synced")
	}
	logrus.Warn("Starting the work for hdfs cluster controller")
	for i := 0 ;i < controllerThreads; i++{
		go wait.Until(c.runWorker, time.Second, ctx.Done())
	}
	return nil
}

// Stop shutdown workQueue
func (c *mysqlController)Stop(stopCh <- chan struct{})error{
	logrus.Info("Stopping the hdfs controller")
	c.queue.ShutDown()
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
		if err := c.operator.Sync(key); err != nil{
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
	//var object metav1.Object
	//var ok bool
	//if object, ok = obj.(metav1.Object); !ok {
	//	tombStone, ok := obj.(cache.DeletedFinalStateUnknown)
	//	if !ok {
	//		runtime.HandleError(fmt.Errorf("error decoding object , wrong type"))
	//		return
	//	}
	//	object, ok = tombStone.Obj.(metav1.Object)
	//	if !ok {
	//		runtime.HandleError(fmt.Errorf("error decoding object tombStone, wrong type"))
	//		return
	//	}
	//}
	key,err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil{
		logrus.Errorf("failed to get key for %v:%v", obj, err)
		return
	}
	c.queue.AddRateLimited(key)
}