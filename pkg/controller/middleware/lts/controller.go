package lts

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
	"l0calh0st.cn/clickpaas-operator/pkg/controller"
	middlewarelister "l0calh0st.cn/clickpaas-operator/pkg/client/listers/middleware/v1alpha1"
	middlewareinformer "l0calh0st.cn/clickpaas-operator/pkg/client/informers/externalversions"
	"l0calh0st.cn/clickpaas-operator/pkg/operator"
	"l0calh0st.cn/clickpaas-operator/pkg/operator/middleware/lts"
	"time"
)

type ltsJobTrackController struct {
	controller.BaseController

	kubeClient kubernetes.Interface
	crdClient crdclient.Interface

	queue workqueue.RateLimitingInterface
	recorder record.EventRecorder
	cacheSyncedList []cache.InformerSynced

	ltsLister middlewarelister.LtsJobTrackerLister

	serviceLister corev1lister.ServiceLister
	deploymentLister appv1lister.DeploymentLister

	operator operator.IOperator
}


func NewLtsJobTrackerController(kubeClient kubernetes.Interface, crdClient crdclient.Interface,
	kubeInformerFactory informers.SharedInformerFactory, midInformerFactory middlewareinformer.SharedInformerFactory,
)*ltsJobTrackController{

	eventBroadCaster := record.NewBroadcaster()
	eventBroadCaster.StartLogging(glog.V(2).Infof)
	eventBroadCaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events(corev1.NamespaceAll)})

	recorder := eventBroadCaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "lts-controller"})

	return newLtsJobTrackerController(kubeClient, crdClient, kubeInformerFactory, midInformerFactory, recorder)

}

func newLtsJobTrackerController(kubeClient kubernetes.Interface, crdClient crdclient.Interface,
	kubeInformerFactory informers.SharedInformerFactory, midInformerFactory middlewareinformer.SharedInformerFactory,
	recorder record.EventRecorder)*ltsJobTrackController{

	controller := &ltsJobTrackController{
		kubeClient:       kubeClient,
		crdClient:        crdClient,
		queue:            workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		recorder:         recorder,
		cacheSyncedList:  []cache.InformerSynced{},
	}

	ltsInformer := midInformerFactory.Middleware().V1alpha1().LtsJobTrackers()
	ltsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.onAdd,
		UpdateFunc: controller.onUpdate,
		DeleteFunc: controller.onDelete,
	})
	controller.cacheSyncedList = append(controller.cacheSyncedList, ltsInformer.Informer().HasSynced)
	controller.ltsLister = ltsInformer.Lister()

	serviceInformer := kubeInformerFactory.Core().V1().Services()
	controller.cacheSyncedList = append(controller.cacheSyncedList, serviceInformer.Informer().HasSynced)
	controller.serviceLister = serviceInformer.Lister()

	deployInformer := kubeInformerFactory.Apps().V1().Deployments()
	controller.cacheSyncedList = append(controller.cacheSyncedList, deployInformer.Informer().HasSynced)
	controller.deploymentLister = deployInformer.Lister()

	controller.operator = lts.NewLtsJobTrackerOperator(kubeClient, controller.ltsLister, controller.deploymentLister, controller.serviceLister)

	return controller
}

func (c *ltsJobTrackController) Start(ctx context.Context, threads int) error {
	if ok := cache.WaitForCacheSync(ctx.Done(),c.cacheSyncedList...); !ok {
		return fmt.Errorf("lts controller, wait for all informer cache synced failed")
	}
	logrus.Infof("lst controller has synced, ready reconcilcing")
	for i := 0 ; i< threads; i++{
		go wait.Until(c.runWorker, 5 * time.Second, ctx.Done())
	}
	<- ctx.Done()
	return ctx.Err()
}

func (c *ltsJobTrackController) Stop(stopCh <-chan struct{}) error {
	panic("implement me")
}


func (c *ltsJobTrackController)runWorker(){
	for c.processNextItem(){}
}
func(c *ltsJobTrackController)processNextItem()bool{
	obj,shutdown := c.queue.Get()
	if shutdown{
		return false
	}
	err := func(obj interface{}) error{
		defer c.queue.Done(obj)
		key, ok := obj.(string)
		if !ok {
			runtime.HandleError(fmt.Errorf("except got string from workqueue, but got %#v", obj))
			return nil
		}
		if err := c.operator.Reconcile(key);err != nil{
			c.queue.AddRateLimited(key)
			return fmt.Errorf("reconciling error %s", err)
		}
		c.queue.Forget(obj)
		return nil
	}(obj)
	if err != nil{
		runtime.HandleError(err)
	}
	return true
}


func(c *ltsJobTrackController)onAdd(obj interface{}){
	lts := obj.(*crdv1alpha1.LtsJobTracker)
	crdv1alpha1.WithDefaultsLtsJobTracker(lts)
	c.recorder.Event(lts,corev1.EventTypeNormal, LtsJobTrackerEventReasonOnAdded, eventInfoMessage(lts, LtsJobTrackerEventReasonOnAdded))
	for _, hook := range c.GetHooks(){
		hook.OnAdd(lts)
	}
	c.enqueue(lts)
}
func(c *ltsJobTrackController)onUpdate(oldObj,newObj interface{}){
	oldLts := oldObj.(*crdv1alpha1.LtsJobTracker)
	newLts := newObj.(*crdv1alpha1.LtsJobTracker)
	if oldLts.ResourceVersion == newLts.ResourceVersion{
		return
	}
	c.recorder.Event(newLts, corev1.EventTypeNormal, LtsJobTrackerEventReasonOnUpdate, eventInfoMessage(newLts, LtsJobTrackerEventReasonOnUpdate))
	for _,hook := range c.GetHooks(){
		hook.OnUpdate(newLts)
	}
	c.enqueue(newLts)
}

func(c *ltsJobTrackController)onDelete(obj interface{}){
	var lts *crdv1alpha1.LtsJobTracker
	switch obj.(type) {
	case *crdv1alpha1.LtsJobTracker:
		lts = obj.(*crdv1alpha1.LtsJobTracker)
	case cache.DeletedFinalStateUnknown:
		deleteObj := obj.(cache.DeletedFinalStateUnknown).Obj
		lts = deleteObj.(*crdv1alpha1.LtsJobTracker)
	}
	if lts == nil{
		return
	}
	c.recorder.Event(lts, corev1.EventTypeNormal, LtsJobTrackerEventReasonOnDelete, eventInfoMessage(lts, LtsJobTrackerEventReasonOnDelete))
	for _,hook := range  c.GetHooks(){
		hook.OnDelete(lts)
	}
	//
}

func (c *ltsJobTrackController)enqueue(obj interface{}){
	key,err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil{
		runtime.HandleError(fmt.Errorf("handle key from workqueue failed, %s", err))
		return
	}
	c.queue.AddRateLimited(key)
}
