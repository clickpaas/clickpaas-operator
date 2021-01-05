package diamond

import (
	"k8s.io/client-go/kubernetes"
	appv1lister "k8s.io/client-go/listers/apps/v1"
	corev1lister "k8s.io/client-go/listers/core/v1"
	crdlister "l0calh0st.cn/clickpaas-operator/pkg/client/listers/middleware/v1alpha1"
)

type diamondOperator struct {
	kubeClient kubernetes.Interface

	diamondLister crdlister.DiamondLister
	serviceLister corev1lister.ServiceLister
	deployLister appv1lister.DeploymentLister
	configmapLister corev1lister.ConfigMapLister
}



func NewDiamondOperator(kubeClient kubernetes.Interface, diamondLister crdlister.DiamondLister, serviceLister corev1lister.ServiceLister,
	deploymentLister appv1lister.DeploymentLister, configMapLister corev1lister.ConfigMapLister)*diamondOperator{
	return &diamondOperator{
		kubeClient:      kubeClient,
		diamondLister:   diamondLister,
		serviceLister:   serviceLister,
		deployLister:    deploymentLister,
		configmapLister: configMapLister,
	}
}


func (d *diamondOperator) Sync(key string) error {
	panic("implement me")
}

func (d *diamondOperator) Healthy() error {
	panic("implement me")
}

