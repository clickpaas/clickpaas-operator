package zookeeper

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	appv1lister "k8s.io/client-go/listers/apps/v1"
	corev1lister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	middlewarelister "l0calh0st.cn/clickpaas-operator/pkg/client/listers/middleware/v1alpha1"
	"l0calh0st.cn/clickpaas-operator/pkg/operator"
	"l0calh0st.cn/clickpaas-operator/pkg/operator/manager"
	kubeutil "l0calh0st.cn/clickpaas-operator/pkg/operator/util/kube"
)

type zookeeperOperator struct {
	kubeClient kubernetes.Interface
	restConfig *rest.Config

	zkLister           middlewarelister.ZookeeperClusterLister
	statefulSetManager operator.StatefulSetManager
	serviceManager     operator.ServiceManager
	configMapManager   operator.ConfigMapManager
	podManager         operator.PodManager
}

func NewZookeeperOperator(kubeClient kubernetes.Interface, restConfig *rest.Config, zkLister middlewarelister.ZookeeperClusterLister,
	statefulSetLister appv1lister.StatefulSetLister, serviceLister corev1lister.ServiceLister,
	configMapLister corev1lister.ConfigMapLister, podLister corev1lister.PodLister) *zookeeperOperator {
	return &zookeeperOperator{
		kubeClient:         kubeClient,
		zkLister:           zkLister,
		restConfig:         restConfig,
		statefulSetManager: manager.NewStatefulSetManager(kubeClient, statefulSetLister),
		serviceManager:     manager.NewServiceManager(kubeClient, serviceLister),
		configMapManager:   manager.NewConfigManager(kubeClient, configMapLister),
		podManager:         manager.NewPodManager(kubeClient, podLister),
	}
}

func (op *zookeeperOperator) Reconcile(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("handler key from workqueue failed %s", err))
		return nil
	}

	zk, err := op.zkLister.ZookeeperClusters(namespace).Get(name)
	if err != nil {
		if k8serr.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("resource zookeeper %s:%s is not existed", name, namespace))
			return nil
		} else {
			return err
		}
	}
	cm, err := op.configMapManager.Get(&configMapResourceEr{zk, newConfigMapForZookeeper})
	if err != nil {
		if !k8serr.IsNotFound(err) {
			return err
		}
		cm, err = op.configMapManager.Create(&configMapResourceEr{zk, newConfigMapForZookeeper})
		if err != nil {
			return err
		}
	}
	_ = cm
	// check service for cluster communicate existed, if not exist ,then create one
	syncSvc, err := op.serviceManager.Get(&serviceResourceEr{zk, newServiceForZookeeperServiceCommunicate})
	if err != nil {
		if !k8serr.IsNotFound(err) {
			return err
		}
		syncSvc, err = op.serviceManager.Create(&serviceResourceEr{zk, newServiceForZookeeperServiceCommunicate})
		if err != nil {
			return err
		}
	}
	_ = syncSvc

	//allWorkNode,err := kubeutil.GetAllWorkNode(op.kubeClient)
	//if err != nil{
	//	return fmt.Errorf("cannot list all worknode, %s",err)
	//}
	//if len(allWorkNode) < int(zk.Spec.Replicas) {
	//
	//}

	allExistedPods, err := op.podManager.List(getLabelForZookeeperCluster(zk))
	if err != nil {
		return err
	}
	err = op.mayCreateOrDeletePodsAccordZkReplicas(zk, allExistedPods)
	if err != nil {
		return err
	}
	// create service for client
	svcCli, err := op.serviceManager.Get(&serviceResourceEr{zk, newServiceForZookeeperClient})
	if err != nil {
		if !k8serr.IsNotFound(err) {
			return err
		}
		svcCli, err = op.serviceManager.Create(&serviceResourceEr{zk, newServiceForZookeeperClient})
		if err != nil {
			return err
		}
	}
	_ = svcCli
	if err := kubeutil.WaitForPodsReady(allExistedPods, 5*time.Second); err != nil {
		return fmt.Errorf("wait all pods ready, %s:%s  %s", zk.GetName(), zk.GetNamespace(), err)
	}
	if len(allExistedPods) != len(kubeutil.FilteredActivePods(allExistedPods)) {
		return fmt.Errorf("second check actived pods %s:%s  failed, may exists some unactived pod", zk.GetName(), zk.GetNamespace())
	}
	//logrus.Infof("all existed pod number is %d", len(allExistedPods))
	//if len(allExistedPods) == 0 {
	//	return fmt.Errorf("all existed pod number is 0")
	//}
	//randPickedPod := allExistedPods[0]
	//time.Sleep(10 * time.Second)
	//if err := doOnceBootStrap(randPickedPod, op.kubeClient, op.restConfig, zk); err != nil {
	//	return fmt.Errorf("bootstrap failed %s", err)
	//}

	return nil
}

func (op *zookeeperOperator) mayCreateOrDeletePodsAccordZkReplicas(cluster *crdv1alpha1.ZookeeperCluster, podList []*corev1.Pod) error {

	zkPodList := generateZookeeperPodList(cluster)

	shouldDelete := getShouldBeDeletedPodList(podList, zkPodList)

	for _, pod := range shouldDelete {
		err := op.podManager.Delete(&podResourceEr{object: pod})
		if err != nil {
			logrus.Errorf(fmt.Sprintf("%s", err.Error()))
			return err
		}
	}

	// install
	shouldInstalled := getShouldInstalledPodList(podList, zkPodList)

	for _, pod := range shouldInstalled {
		id, _ := strconv.Atoi(pod.id)
		_, err := op.podManager.Create(&podResourceEr{cluster, id})
		if err != nil {
			logrus.Errorf(fmt.Sprintf("%s", err.Error()))
			return err
		}
	}
	return nil
}

func (op *zookeeperOperator) syncStatefulSet() error {
	return nil
}

func (op *zookeeperOperator) Healthy() error {
	// 检查server client是否存在
	return nil
}
