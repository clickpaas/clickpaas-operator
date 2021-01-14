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
	configMapManager operator.ConfigMapManager
}


func NewRocketmqOperator(kubeClient kubernetes.Interface, rocketmqLister middlewarelister.RocketmqLister,
	deploymentLister appv1lister.DeploymentLister, statefulSetLister appv1lister.StatefulSetLister,
	serviceLister corev1lister.ServiceLister, configMapLister corev1lister.ConfigMapLister)operator.IOperator{
	return &rocketmqOperator{
		kubeClient:         kubeClient,
		rocketmqLister:     rocketmqLister,
		deploymentManager:  manager.NewDeploymentManager(kubeClient, deploymentLister),
		statefulSetManager: manager.NewStatefulSetManager(kubeClient, statefulSetLister),
		serviceManager:     manager.NewServiceManager(kubeClient, serviceLister),
		configMapManager: manager.NewConfigManager(kubeClient, configMapLister),
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
	// check nameserver application is existed,
	nsvcDeploy,err := op.deploymentManager.Get(&deploymentResourceEr{rocket})
	if err != nil{
		if k8serr.IsNotFound(err){
			// if nameserver is not existed, then create an new one
			nsvcDeploy, err = op.deploymentManager.Create(&deploymentResourceEr{rocket})
		}
		if err != nil{
			return err
		}
	}
	_ = nsvcDeploy
	// check nameserver service is existed, if not existed, then create new one
	nSvc,err := op.serviceManager.Get(&serviceResourceEr{rocket, newServiceForRocketmqNameServer})
	if err != nil{
		if k8serr.IsNotFound(err){
			nSvc, err = op.serviceManager.Create(&serviceResourceEr{rocket, newServiceForRocketmqNameServer})
		}
		if err != nil{
			return err
		}
	}
	_ = nSvc
	//
	cm,err := op.configMapManager.Get(&configMapEr{rocket})
	if err != nil{
		if k8serr.IsNotFound(err){
			cm,err = op.configMapManager.Create(&configMapEr{rocket})
		}
		if err != nil{
			return err
		}
	}
	_ = cm

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

	return nil

}

func (op *rocketmqOperator) Healthy() error {
	panic("implement me")
}
