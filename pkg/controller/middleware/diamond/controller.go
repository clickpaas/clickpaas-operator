package diamond

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
	"l0calh0st.cn/clickpaas-operator/pkg/operator/middleware/diamond"
	"time"
)

type diamondController struct {
	controller.BaseController
	kubeClient kubernetes.Interface
	crdClient crdclient.Interface

	diamondLister crdlister.DiamondLister
	deployLister appv1lister.DeploymentLister
	serviceLister corev1lister.ServiceLister
	configMapLister corev1lister.ConfigMapLister

	queue workqueue.RateLimitingInterface
	cacheSyncedList []cache.InformerSynced
	recorder record.EventRecorder

	operator operator.IOperator
}


func NewDiamondController(
	kubeClient kubernetes.Interface,
	crdClient crdclient.Interface,
	kubeInformerFactory informers.SharedInformerFactory,
	crdInformerFactory crdinformer.SharedInformerFactory,
	)*diamondController{
	eventBroadCaster := record.NewBroadcaster()
	eventBroadCaster.StartLogging(glog.V(2).Infof)
	eventBroadCaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events(metav1.NamespaceAll)})
	recorder := eventBroadCaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "diamond-controller"})
	return newDiamondController(kubeClient, crdClient, kubeInformerFactory, crdInformerFactory, recorder)

}

func newDiamondController(kubeClient kubernetes.Interface,
	crdClient crdclient.Interface, kubeInformerFactory informers.SharedInformerFactory,
	crdInformerFactory crdinformer.SharedInformerFactory, recorder record.EventRecorder)*diamondController{

	controller := &diamondController{
		kubeClient:      kubeClient,
		crdClient:       crdClient,
		queue:           workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		cacheSyncedList: []cache.InformerSynced{},
		recorder:        recorder,
	}

	diamondInformer := crdInformerFactory.Middleware().V1alpha1().Diamonds()
	diamondInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.onAdd,
		UpdateFunc: controller.onUpdate,
		DeleteFunc: controller.onDelete,
	})
	controller.cacheSyncedList = append(controller.cacheSyncedList, diamondInformer.Informer().HasSynced)
	controller.diamondLister = diamondInformer.Lister()

	deployInformer := kubeInformerFactory.Apps().V1().Deployments()
	controller.deployLister = deployInformer.Lister()
	controller.cacheSyncedList = append(controller.cacheSyncedList, deployInformer.Informer().HasSynced)

	serviceInformer := kubeInformerFactory.Core().V1().Services()
	controller.serviceLister = serviceInformer.Lister()
	controller.cacheSyncedList = append(controller.cacheSyncedList, serviceInformer.Informer().HasSynced)

	configMapInformer := kubeInformerFactory.Core().V1().ConfigMaps()
	controller.configMapLister = configMapInformer.Lister()
	controller.cacheSyncedList = append(controller.cacheSyncedList, configMapInformer.Informer().HasSynced)

	controller.operator = diamond.NewDiamondOperator(kubeClient, controller.diamondLister, controller.serviceLister, controller.deployLister, controller.configMapLister)
	return controller
}

func (d *diamondController)onAdd(obj interface{}){
	diamond := obj.(*crdv1alpha1.Diamond)
	crdv1alpha1.WithDefaultsDiamond(diamond)
	logrus.Info("Diamond  %v was added,enqueue for next handling", diamond.GetName())
	d.recorder.Event(diamond, corev1.EventTypeNormal, DiamondEventReasonOnAdd, fmt.Sprintf("%v created", diamond.GetName()))
	for _,hook := range d.GetHooks(){
		hook.OnAdd(diamond)
	}
	d.enqueue(diamond)
}

func(d *diamondController)onDelete(obj interface{}){
	var diamond *crdv1alpha1.Diamond
	switch obj.(type) {
	case *crdv1alpha1.Diamond:
		diamond = obj.(*crdv1alpha1.Diamond)
	case cache.DeletedFinalStateUnknown:
		deleteObj := obj.(cache.DeletedFinalStateUnknown).Obj
		diamond = deleteObj.(*crdv1alpha1.Diamond)
	}
	if diamond != nil{
	}
	d.recorder.Event(diamond, corev1.EventTypeNormal, DiamondEventReasonOnDelete, fmt.Sprintf("%v deleted", diamond.GetName()))

	for _, hook := range d.GetHooks(){
		hook.OnDelete(diamond)
	}
	// todo if need for next handler this object , or do nothing
}

func (d *diamondController)onUpdate(oldObj,newObj interface{}){
	oldDiamond := oldObj.(*crdv1alpha1.Diamond)
	newDiamond := newObj.(*crdv1alpha1.Diamond)
	d.recorder.Event(newDiamond, corev1.EventTypeNormal, DiamondEventReasonOnUpdate, fmt.Sprintf("%v updated", newDiamond.GetName()))
	// if resource  has no change
	if oldDiamond.ResourceVersion == newDiamond.ResourceVersion{
		// resource has no change, do nothing and return
		return
	}
	//if !equality.Semantic.DeepEqual(oldDiamond, newDiamond){
	//	// if semantic changed, then delete old ,and crate new one
	//	err := d.crdClient.MiddlewareV1alpha1().Diamonds(oldDiamond.GetNamespace()).Delete(context.TODO(), oldDiamond.GetName(), metav1.DeleteOptions{})
	//	if err != nil{
	//	}
	//	newDiamond,err= d.crdClient.MiddlewareV1alpha1().Diamonds(newDiamond.GetNamespace()).Create(context.TODO(), newDiamond, metav1.CreateOptions{})
	//}
	for _,hook := range d.GetHooks(){
		hook.OnUpdate(newDiamond)
	}
	d.enqueue(newDiamond)
}


func (d *diamondController) Start(ctx context.Context, threads int) error {
	logrus.Warn("ready start diamond Controller..waiting all informer cache hasSynced")
	if ok := cache.WaitForCacheSync(ctx.Done(), d.cacheSyncedList...); !ok {
		return fmt.Errorf("wait all informer cache synced failed")
	}
	logrus.Infof("Diamon Controller has stared, ready to reclioning")
	for i := 0; i<threads; i++{
		go wait.Until(d.runWorker, 5 * time.Second, ctx.Done())
	}
	<- ctx.Done()
	return ctx.Err()
}

func(d *diamondController)runWorker(){
	defer runtime.HandleCrash()
	for d.processNextItem(){}
}

func(d *diamondController)processNextItem()bool{
	obj,shutdown := d.queue.Get()
	if shutdown{
		return false
	}
	err := func(obj interface{}) error{
		defer d.queue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			d.queue.Forget(obj)
			runtime.HandleError(fmt.Errorf("except got string in workqueue, but got %#v", obj))
			return nil
		}
		if err := d.operator.Reconcile(key); err != nil{
			d.queue.AddRateLimited(key)
			return fmt.Errorf("error syncing diamond %v:%v, requeue", key, err)
		}
		d.queue.Forget(obj)
		logrus.Infof("successfully synced %s", key)
		return nil
	}(obj)

	if err != nil{
		runtime.HandleError(err)
		return true
	}
	return true
}

func (d *diamondController) Stop(stopCh <-chan struct{}) error {
	d.queue.ShutDown()
	return nil
}



func (d *diamondController)enqueue(obj interface{}){
	key,err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil{
		logrus.Errorf("failed to get key for %v:%v", obj, err)
		return
	}
	d.queue.AddRateLimited(key)
}

