package diamond

import (
	"k8s.io/client-go/kubernetes"
	corev1lister "k8s.io/client-go/listers/core/v1"
)




type configMapManager struct {
	kubeClient kubernetes.Interface
	configMapLister corev1lister.ConfigMapLister
}

func NewConfigMapManager(kubeClient kubernetes.Interface, configMapLister corev1lister.ConfigMapLister)*configMapManager{
	return &configMapManager{
		kubeClient:      kubeClient,
		configMapLister: configMapLister,
	}
}