package diamond

import (
	"fmt"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	appv1lister "k8s.io/client-go/listers/apps/v1"
	corev1lister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	crdlister "l0calh0st.cn/clickpaas-operator/pkg/client/listers/middleware/v1alpha1"
)

type diamondOperator struct {
	kubeClient kubernetes.Interface

	diamondLister crdlister.DiamondLister
	serviceManager *serviceManager
	configMapManager *configMapManager
	deploymentManager *deploymentManager
}



func NewDiamondOperator(kubeClient kubernetes.Interface, diamondLister crdlister.DiamondLister, serviceLister corev1lister.ServiceLister,
	deploymentLister appv1lister.DeploymentLister, configMapLister corev1lister.ConfigMapLister)*diamondOperator{
	op := &diamondOperator{
		kubeClient:      kubeClient,
		diamondLister:   diamondLister,

	}
	op.configMapManager = NewConfigMapManager(kubeClient, configMapLister)
	op.serviceManager = NewServiceManager(kubeClient, serviceLister)
	op.deploymentManager = NewDeploymentManager(kubeClient, deploymentLister)
	return op
}


func (d *diamondOperator) Sync(key string) error {
	namespace,name,err := cache.SplitMetaNamespaceKey(key)
	if err != nil{
		runtime.HandleError(fmt.Errorf("split key error '%v':%v", key, err))
		return nil
	}
	diamond,err := d.diamondLister.Diamonds(namespace).Get(name)
	if err != nil{
		if k8serr.IsNotFound(err){
			runtime.HandleError(fmt.Errorf("diamond is not existe"))
			return nil
		}
		return nil
	}
	// check deployment is existed
	deploy,err := d.deploymentManager.Get(diamond)
	if err != nil{
		if k8serr.IsNotFound(err){
			deploy,err = d.deploymentManager.Create(diamond)
			if err != nil{
				return err
			}
		} else {
			return err
		}
	}
	if *deploy.Spec.Replicas != diamond.Spec.Replicas{
		deploy.Spec.Replicas = &diamond.Spec.Replicas
		deploy,err = d.deploymentManager.Update(deploy)
		if err != nil{
			return err
		}
	}
	// check service
	service,err := d.serviceManager.Get(diamond)
	if err != nil{
		if k8serr.IsNotFound(err){
			service,err = d.serviceManager.Create(diamond)
			if err != nil{
				return err
			}
		} else {
			return err
		}
	}
	_ = service
	return nil
}

func (d *diamondOperator) Healthy() error {
	return nil
}

