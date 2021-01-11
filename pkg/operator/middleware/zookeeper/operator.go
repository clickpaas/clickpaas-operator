package zookeeper

import (
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	appv1lister "k8s.io/client-go/listers/apps/v1"
	corev1lister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	middlewarelister "l0calh0st.cn/clickpaas-operator/pkg/client/listers/middleware/v1alpha1"
	"l0calh0st.cn/clickpaas-operator/pkg/operator"
	"l0calh0st.cn/clickpaas-operator/pkg/operator/manager"
	"strconv"
)

type zookeeperOperator struct {
	kubeClient kubernetes.Interface

	zkLister middlewarelister.ZookeeperClusterLister
	statefulSetManager operator.StatefulSetManager
	serviceManager operator.ServiceManager
	configMapManager operator.ConfigMapManager
	podManager operator.PodManager

}


func NewZookeeperOperator(kubeClient kubernetes.Interface, zkLister middlewarelister.ZookeeperClusterLister,
	statefulSetLister appv1lister.StatefulSetLister, serviceLister corev1lister.ServiceLister,
	configMapLister corev1lister.ConfigMapLister, podLister corev1lister.PodLister)*zookeeperOperator{
	return &zookeeperOperator{
		kubeClient:         kubeClient,
		zkLister:           zkLister,
		statefulSetManager: manager.NewStatefulSetManager(kubeClient, statefulSetLister),
		serviceManager:     manager.NewServiceManager(kubeClient, serviceLister),
		configMapManager: manager.NewConfigManager(kubeClient,configMapLister),
		podManager: manager.NewPodManager(kubeClient, podLister),
	}
}


func (op *zookeeperOperator) Reconcile(key string) error {
	namespace,name,err := cache.SplitMetaNamespaceKey(key)
	if err != nil{
		runtime.HandleError(fmt.Errorf("handler key from workqueue failed %s", err))
		return nil
	}

	zk,err := op.zkLister.ZookeeperClusters(namespace).Get(name)
	if err != nil{
		if k8serr.IsNotFound(err){
			runtime.HandleError(fmt.Errorf("resource zookeeper %s:%s is not existed", name, namespace))
			return nil
		}else {
			return err
		}
	}
	logrus.Infof("raday handler config")
	cm,err := op.configMapManager.Get(&configMapResourceEr{zk, newConfigMapForZookeeper})
	if err != nil{
		if !k8serr.IsNotFound(err){
			return err
		}
		cm,err = op.configMapManager.Create(&configMapResourceEr{zk, newConfigMapForZookeeper})
		if err != nil{
			return err
		}
	}
	_ = cm
	logrus.Infof("ready to handler service")
	// check service for cluster communicate existed, if not exist ,then create one
	syncSvc,err := op.serviceManager.Get(&serviceResourceEr{zk,newServiceForZookeeperServiceCommunicate})
	if err != nil {
		if !k8serr.IsNotFound(err){
			return err
		}
		syncSvc,err = op.serviceManager.Create(&serviceResourceEr{zk, newServiceForZookeeperServiceCommunicate})
		if err != nil{
			return err
		}
	}
	_ = syncSvc

	allExistedPods,err := op.podManager.List(getLabelForZookeeperCluster(zk))
	logrus.Infof("ready handler pod ,now pod number is %d", len(allExistedPods))
	if err != nil{
		return err
	}
	err = op.mayCreateOrDeletePodsAccordZkReplicas(zk, allExistedPods)
	if err != nil{
		return err
	}
	//
	return nil

}


func(op *zookeeperOperator)mayCreateOrDeletePodsAccordZkReplicas(cluster *crdv1alpha1.ZookeeperCluster,podList []*corev1.Pod)error{

	zkPodList := generateZookeeperPodList(cluster)

	shouldDelete := getShouldBeDeletedPodList(podList, zkPodList)

	for _, pod := range shouldDelete{
		err := op.podManager.Delete(&podResourceEr{object: pod})
		if err != nil{
			logrus.Errorf(fmt.Sprintf("%s", err.Error()))
			return err
		}
	}

	// install
	shouldInstalled := getShouldInstalledPodList(podList, zkPodList)

	for _, pod := range shouldInstalled{
		id,_ := strconv.Atoi(pod.id)
		_,err := op.podManager.Create(&podResourceEr{cluster, id})
		if err != nil{
			logrus.Errorf(fmt.Sprintf("%s", err.Error()))
			return err
		}
	}
	logrus.Infof("may create or delete %d  %d  %d", len(podList), len(shouldDelete), len(shouldInstalled))
	return nil
}




func(op *zookeeperOperator)syncStatefulSet()error{
	return nil
}

func (op *zookeeperOperator) Healthy() error {
	// 检查server client是否存在
	return nil
}