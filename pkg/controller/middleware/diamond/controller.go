package diamond

import (
	"context"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appv1lister "k8s.io/client-go/listers/apps/v1"
	corev1lister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	crdclient "l0calh0st.cn/clickpaas-operator/pkg/client/clientset/versioned"
	"l0calh0st.cn/clickpaas-operator/pkg/client/clientset/versioned/scheme"
	crdinformer "l0calh0st.cn/clickpaas-operator/pkg/client/informers/externalversions"
	crdlister "l0calh0st.cn/clickpaas-operator/pkg/client/listers/middleware/v1alpha1"
	"l0calh0st.cn/clickpaas-operator/pkg/controller"
	"l0calh0st.cn/clickpaas-operator/pkg/operator"
	"l0calh0st.cn/clickpaas-operator/pkg/operator/diamond"
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

}

func(d *diamondController)onDelete(obj interface{}){

}

func (d *diamondController)onUpdate(oldObj,newObj interface{}){

}


func (d *diamondController) Start(ctx context.Context, threads int) error {
	panic("implement me")
}

func (d *diamondController) Stop(stopCh <-chan struct{}) error {
	panic("implement me")
}

func (d *diamondController) AddHook(hook controller.IHook) error {
	panic("implement me")
}

func (d *diamondController) RemoveHook(hook controller.IHook) error {
	panic("implement me")
}


