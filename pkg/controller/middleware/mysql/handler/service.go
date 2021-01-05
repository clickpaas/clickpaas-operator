package handler

import (
	"k8s.io/client-go/kubernetes"
	corev1lister "k8s.io/client-go/listers/core/v1"
)

type serviceEventHandler struct {
	kubeClient kubernetes.Interface
	serviceLister corev1lister.ServiceLister
}

func NewServiceManager(kubeClient kubernetes.Interface, svcLister corev1lister.ServiceLister)*serviceEventHandler{
	return &serviceEventHandler{
		kubeClient:    kubeClient,
		serviceLister: svcLister,
	}
}


func OnAddEventHandler(obj interface{}){

}

func OnUpdateEventHandler(oldObj,newObj interface{}){

}

func OnDeleteEventHandler(obj interface{}){

}
