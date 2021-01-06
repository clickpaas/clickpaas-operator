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
	Create(interface{}, func(interface{})(*appv1.Deployment,error))(*appv1.Deployment,error)
	Update(interface{}, func(interface{})(*appv1.Deployment,error))(*appv1.Deployment,error)
	Delete(interface{}, func(interface{})(*appv1.Deployment,error))error
	Get(interface{}, func(interface{})(*appv1.Deployment,error))(*appv1.Deployment,error)
}


type ConfigMapManager interface {
	Create(interface{}, func(interface{})(*corev1.ConfigMap,error))(*corev1.ConfigMap,error)
	Update(interface{}, func(interface{})(*corev1.ConfigMap,error))(*corev1.ConfigMap,error)
	Delete(interface{}, func(interface{})(*corev1.ConfigMap,error))error
	Get(interface{}, func(interface{})(*corev1.ConfigMap,error))(*corev1.ConfigMap,error)
}

type StatefulSetManager interface {
	Create(interface{}, func(interface{})(*appv1.StatefulSet,error))(*appv1.StatefulSet, error)
	Update(interface{}, func(interface{})(*appv1.StatefulSet,error))(*appv1.StatefulSet, error)
	Delete(interface{}, func(interface{})(*appv1.StatefulSet,error))error
	Get(interface{}, func(interface{})(*appv1.StatefulSet,error))(*appv1.StatefulSet, error)
}


type ServiceManager interface {
	Create(interface{}, func(interface{})(*corev1.Service,error))(*corev1.Service, error)
	Update(interface{}, func(interface{})(*corev1.Service,error))(*corev1.Service, error)
	Delete(interface{}, func(interface{})(*corev1.Service,error))error
	Get(interface{}, func(interface{})(*corev1.Service,error))(*corev1.Service, error)
}