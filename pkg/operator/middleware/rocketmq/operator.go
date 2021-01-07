package rocketmq

import (
	"fmt"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	appv1lister "k8s.io/client-go/listers/apps/v1"
	corev1lister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	middlewarelister "l0calh0st.cn/clickpaas-operator/pkg/client/listers/middleware/v1alpha1"
	"l0calh0st.cn/clickpaas-operator/pkg/operator"
	"l0calh0st.cn/clickpaas-operator/pkg/operator/manager"
)

type rocketmqOperator struct {
	kubeClient kubernetes.Interface

	rocketmqLister middlewarelister.RocketmqLister

	deploymentManager operator.DeploymentManager
	statefulSetManager operator.StatefulSetManager
	serviceManager operator.ServiceManager
}


func NewRocketmqOperator(kubeClient kubernetes.Interface, rocketmqLister middlewarelister.RocketmqLister,
	deploymentLister appv1lister.DeploymentLister, statefulSetLister appv1lister.StatefulSetLister,
	serviceLister corev1lister.ServiceLister)operator.IOperator{
	return &rocketmqOperator{
		kubeClient:         kubeClient,
		rocketmqLister:     rocketmqLister,
		deploymentManager:  manager.NewDeploymentManager(kubeClient, deploymentLister),
		statefulSetManager: manager.NewStatefulSetManager(kubeClient, statefulSetLister),
		serviceManager:     manager.NewServiceManager(kubeClient, serviceLister),
	}
}


func (op *rocketmqOperator) Reconcile(key string) error {
	namespace,name,err := cache.SplitMetaNamespaceKey(key)
	if err != nil{
		runtime.HandleError(fmt.Errorf("splitMetanamespacekey error, %s", err))
		return nil
	}
	rocket,err := op.rocketmqLister.Rocketmqs(namespace).Get(name)
	if err != nil{
		if k8serr.IsNotFound(err){
			runtime.HandleError(fmt.Errorf("rocketmq is not exist '%s:%s'", name, namespace))
			return nil
		}
		return err
	}
	// start nameserver
	dp,err := op.deploymentManager.Get(&deploymentResourceEr{rocket})
	if err != nil{
		if k8serr.IsNotFound(err){
			dp,err = op.deploymentManager.Create(&deploymentResourceEr{rocket})
		}
		if err != nil{
			return err
		}
	}
	_ = dp
	// check nameserver service
	nsSvc ,err := op.serviceManager.Get(&serviceResourceErForNameServer{rocket})
	if err != nil{
		if k8serr.IsNotFound(err){
			nsSvc,err = op.serviceManager.Create(&serviceResourceErForNameServer{rocket})
		}
		if err != nil{
			return err
		}
	}
	_ = nsSvc
	// start statefulset
	ss,err := op.statefulSetManager.Get(&statefulSetResourceEr{rocket})
	if err != nil{
		if k8serr.IsNotFound(err){
			ss,err = op.statefulSetManager.Create(&statefulSetResourceEr{rocket})
		}
		if err != nil{
			return err
		}
	}
	_ = ss
	// start rocketmqa service
	rkSvc,err := op.serviceManager.Get(&serviceResourceErForRocketmq{rocket})
	if err != nil{
		if k8serr.IsNotFound(err){
			rkSvc, err = op.serviceManager.Create(&serviceResourceErForRocketmq{rocket})
		}
		if err != nil{
			return err
		}
	}
	_ = rkSvc
	return nil

}

func (op *rocketmqOperator) Healthy() error {
	panic("implement me")
}
