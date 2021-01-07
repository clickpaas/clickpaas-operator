package lts

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

type ltsOperator struct {
	kubeClient kubernetes.Interface

	ltsLister middlewarelister.LtsJobTrackerLister

	deploymentManager operator.DeploymentManager
	serviceManager operator.ServiceManager
}



func NewLtsJobTrackerOperator(kubeClient kubernetes.Interface, ltsLister middlewarelister.LtsJobTrackerLister,
	deployLister appv1lister.DeploymentLister, serviceLister corev1lister.ServiceLister)operator.IOperator{
	op := &ltsOperator{
		kubeClient:        kubeClient,
		ltsLister:         ltsLister,
	}
	op.deploymentManager = manager.NewDeploymentManager(kubeClient, deployLister)
	op.serviceManager = manager.NewServiceManager(kubeClient, serviceLister)
	return op
}
func (op *ltsOperator) Reconcile(key string) error {
	namespace,name,err := cache.SplitMetaNamespaceKey(key)
	if err != nil{
		runtime.HandleError(fmt.Errorf(fmt.Sprintf("split key failed %s", err)))
		return nil
	}
	lts,err := op.ltsLister.LtsJobTrackers(namespace).Get(name)
	if err != nil{
		if k8serr.IsNotFound(err){
			runtime.HandleError(fmt.Errorf("lts %s:%s is not existed", namespace, name))
			return nil
		}
		return err
	}
	dp,err := op.deploymentManager.Get(&deploymentResourceEr{lts})
	if err != nil{
		if k8serr.IsNotFound(err){
			dp,err = op.deploymentManager.Create(&deploymentResourceEr{lts})
		}
		if err != nil{
			return err
		}
	}
	_ = dp
	return nil
}

func (op *ltsOperator) Healthy() error {
	panic("implement me")
}