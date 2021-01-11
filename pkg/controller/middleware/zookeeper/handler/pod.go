package handler

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"

	middlewarelister "l0calh0st.cn/clickpaas-operator/pkg/client/listers/middleware/v1alpha1"
)

type podEventHandler struct {
	zkLister middlewarelister.ZookeeperClusterLister
	enqueueFunc func(interface{})
}

func NewPodEventHandler(zkLister middlewarelister.ZookeeperClusterLister, enqueueFunc func(interface{}))*podEventHandler{
	return &podEventHandler{
		zkLister:    zkLister,
		enqueueFunc: enqueueFunc,
	}
}


func(h *podEventHandler)OnAdd(obj interface{}){
	pod,ok := obj.(*corev1.Pod)
	if !ok {
		return
	}
	h.mayEnqueueZkWhenPodChanged(pod)
}

func(h *podEventHandler)OnUpdate(oldObj, newObj interface{}){
	oldPod,ok := oldObj.(*corev1.Pod)
	if !ok {
		return
	}
	newPod,ok := newObj.(*corev1.Pod)
	if !ok {
		return
	}
	if oldPod.ResourceVersion == newPod.ResourceVersion{
		return
	}

	h.mayEnqueueZkWhenPodChanged(newPod)
}

func(h *podEventHandler)OnDelete(obj interface{}){
	pod,ok := obj.(*corev1.Pod)
	if !ok {
		return
	}
	h.mayEnqueueZkWhenPodChanged(pod)
}

func(h *podEventHandler)mayEnqueueZkWhenPodChanged(pod *corev1.Pod){
	if len(pod.OwnerReferences) == 0{
		return
	}
	ownerReference := pod.OwnerReferences[0]
	zk,err := h.zkLister.ZookeeperClusters(pod.GetNamespace()).Get(ownerReference.Name)
	if err != nil{
		return
	}
	for _, ownerference := range pod.OwnerReferences{
		if ownerference.Kind != crdv1alpha1.ZookeeperKind || !metav1.IsControlledBy(pod, zk){

			return
		}
	}
	h.enqueueFunc(zk)
}

