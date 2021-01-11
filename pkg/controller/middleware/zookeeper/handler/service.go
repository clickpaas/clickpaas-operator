package handler

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	middlewarelister "l0calh0st.cn/clickpaas-operator/pkg/client/listers/middleware/v1alpha1"
)

type serviceEventHandler struct {
	zkLister middlewarelister.ZookeeperClusterLister
	enqueueFunc func(interface{})
}

func NewServiceEventHandler(zkLister middlewarelister.ZookeeperClusterLister, enqueueFunc func(interface{}))*serviceEventHandler{
	return &serviceEventHandler{
		zkLister:    zkLister,
		enqueueFunc: enqueueFunc,
	}
}


func(h *serviceEventHandler)OnAdd(obj interface{}){
	svc,ok := obj.(*corev1.Service)
	if !ok {
		return
	}
	h.mayEnqueueZkWhenServiceChanged(svc)
}

func(h *serviceEventHandler)OnUpdate(oldObj, newObj interface{}){
	oldSvc,ok := oldObj.(*corev1.Service)
	if !ok {
		return
	}
	newSvc,ok := newObj.(*corev1.Service)
	if !ok {
		return
	}
	if oldSvc.ResourceVersion == newSvc.ResourceVersion{
		return
	}
	h.mayEnqueueZkWhenServiceChanged(newSvc)
}

func(h *serviceEventHandler)OnDelete(obj interface{}){
	svc,ok := obj.(*corev1.Service)
	if !ok {
		return
	}
	h.mayEnqueueZkWhenServiceChanged(svc)
}

func(h *serviceEventHandler)mayEnqueueZkWhenServiceChanged(svc *corev1.Service){
	if len(svc.OwnerReferences) == 0{
		return
	}
	ownerReference := svc.OwnerReferences[0]
	zk,err := h.zkLister.ZookeeperClusters(svc.GetNamespace()).Get(ownerReference.Name)
	if err != nil{
		return
	}
	for _, ownerference := range svc.OwnerReferences{
		if ownerference.Kind != crdv1alpha1.ZookeeperKind || !metav1.IsControlledBy(svc, zk){
			return
		}
	}
	h.enqueueFunc(zk)
}