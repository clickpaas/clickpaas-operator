package operator

import (
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)
type PodResourceEr interface {
	PodResourceEr(...interface{})(*corev1.Pod,error)
}
type DeploymentResourceEr interface {
	DeploymentResourceEr(...interface{})(*appv1.Deployment, error)
}
type ServiceResourceEr interface {
	ServiceResourceEr(...interface{})(*corev1.Service,error)
}
type ConfigMapResourceEr interface {
	ConfigMapResourceEr(...interface{})(*corev1.ConfigMap,error)
}
type StatefulSetResourceEr interface {
	StatefulSetResourceEr(...interface{})(*appv1.StatefulSet,error)
}
