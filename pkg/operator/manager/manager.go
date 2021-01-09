package manager

import (
	"context"
	"fmt"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appv1lister "k8s.io/client-go/listers/apps/v1"
	corev1lister "k8s.io/client-go/listers/core/v1"
)

type kubeResManager struct {
	kubeClient kubernetes.Interface
	podLister corev1lister.PodLister
	serviceLister corev1lister.ServiceLister
	statefulSetLister appv1lister.StatefulSetLister
	configMapLister corev1lister.ConfigMapLister
	deploymentLister appv1lister.DeploymentLister
}

func NewKubeResourceManager(kubeClient kubernetes.Interface, podLister corev1lister.PodLister, serviceLister corev1lister.ServiceLister,
	statefulSetLister appv1lister.StatefulSetLister, configMapLister corev1lister.ConfigMapLister, deploymentLister appv1lister.DeploymentLister,
	)*kubeResManager{
	manager := &kubeResManager{
		kubeClient:        kubeClient,
		podLister:         podLister,
		serviceLister:     serviceLister,
		statefulSetLister: statefulSetLister,
		configMapLister:   configMapLister,
		deploymentLister:  deploymentLister,
	}
	return manager
}


func(m *kubeResManager)Create(obj interface{})(interface{},error){
	switch obj.(type) {
	case *corev1.Pod:
		pod := obj.(*corev1.Pod)
		return m.kubeClient.CoreV1().Pods(pod.GetNamespace()).Create(context.TODO(), pod, metav1.CreateOptions{})
	case *corev1.ConfigMap:
		cm := obj.(*corev1.ConfigMap)
		return m.kubeClient.CoreV1().ConfigMaps(cm.GetNamespace()).Create(context.TODO(), cm, metav1.CreateOptions{})
	case *appv1.Deployment:
		dp := obj.(*appv1.Deployment)
		return m.kubeClient.AppsV1().Deployments(dp.GetNamespace()).Create(context.TODO(), dp, metav1.CreateOptions{})
	case *appv1.StatefulSet:
		ss := obj.(*appv1.StatefulSet)
		return m.kubeClient.AppsV1().StatefulSets(ss.GetNamespace()).Create(context.TODO(), ss, metav1.CreateOptions{})

	}
	return nil, fmt.Errorf("unknown type object  %#v", obj)
}


