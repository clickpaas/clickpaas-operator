package operator

import (
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

type Manager struct {
	KubeClient kubernetes.Interface
}


type KubeResManager interface {
	Create(interface{})(interface{},error)
	Get(interface{})(interface{}, error)
	Update(interface{})(interface{},error)
	Delete(interface{})error
	List(set labels.Set)([]interface{}, error)
}

type PodManager interface {
	Create(PodResourceEr)(*corev1.Pod,error)
	Get(PodResourceEr)(*corev1.Pod,error)
	Update(PodResourceEr)(*corev1.Pod,error)
	Delete(PodResourceEr)error
	List(labels.Set)([]*corev1.Pod, error)
}


type DeploymentManager interface {
	Create(DeploymentResourceEr)(*appv1.Deployment,error)
	Update(DeploymentResourceEr)(*appv1.Deployment,error)
	Delete(DeploymentResourceEr)error
	Get(DeploymentResourceEr)(*appv1.Deployment,error)
	List(labels.Set)([]*appv1.Deployment, error)
}


type ConfigMapManager interface {
	Create(ConfigMapResourceEr)(*corev1.ConfigMap,error)
	Update(ConfigMapResourceEr)(*corev1.ConfigMap,error)
	Delete(ConfigMapResourceEr)error
	Get(ConfigMapResourceEr)(*corev1.ConfigMap,error)
	List(labels.Set)([]*corev1.ConfigMap, error)
}


type StatefulSetManager interface {
	Create(StatefulSetResourceEr)(*appv1.StatefulSet, error)
	Update(StatefulSetResourceEr)(*appv1.StatefulSet, error)
	Delete(StatefulSetResourceEr)error
	Get(StatefulSetResourceEr)(*appv1.StatefulSet, error)
	List(labels.Set)([]*appv1.StatefulSet, error)
}


type ServiceManager interface {
	Create(ServiceResourceEr)(*corev1.Service, error)
	Update(ServiceResourceEr)(*corev1.Service, error)
	Delete(ServiceResourceEr)error
	Get(ServiceResourceEr)(*corev1.Service, error)
	List(labels.Set)([]*corev1.Service, error)
}