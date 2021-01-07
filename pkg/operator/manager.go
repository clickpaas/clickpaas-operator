package operator

import (
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)


type Manager struct {
	KubeClient kubernetes.Interface
}


type DeploymentManager interface {
	Create(DeploymentResourceEr)(*appv1.Deployment,error)
	Update(DeploymentResourceEr)(*appv1.Deployment,error)
	Delete(DeploymentResourceEr)error
	Get(DeploymentResourceEr)(*appv1.Deployment,error)
}


type ConfigMapManager interface {
	Create(ConfigMapResourceEr)(*corev1.ConfigMap,error)
	Update(ConfigMapResourceEr)(*corev1.ConfigMap,error)
	Delete(ConfigMapResourceEr)error
	Get(ConfigMapResourceEr)(*corev1.ConfigMap,error)
}


type StatefulSetManager interface {
	Create(StatefulSetResourceEr)(*appv1.StatefulSet, error)
	Update(StatefulSetResourceEr)(*appv1.StatefulSet, error)
	Delete(StatefulSetResourceEr)error
	Get(StatefulSetResourceEr)(*appv1.StatefulSet, error)
}


type ServiceManager interface {
	Create(ServiceResourceEr)(*corev1.Service, error)
	Update(ServiceResourceEr)(*corev1.Service, error)
	Delete(ServiceResourceEr)error
	Get(ServiceResourceEr)(*corev1.Service, error)
}