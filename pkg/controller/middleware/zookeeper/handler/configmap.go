package handler

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"

	middlewarelister "l0calh0st.cn/clickpaas-operator/pkg/client/listers/middleware/v1alpha1"
)

type configMapEventHandler struct {
	zkLister middlewarelister.ZookeeperClusterLister
	enqueueFunc func(interface{})
}

func NewConfigMapEventHandler(zkLister middlewarelister.ZookeeperClusterLister, enqueueFunc func(interface{}))*configMapEventHandler{
	return &configMapEventHandler{
		zkLister:    zkLister,
		enqueueFunc: enqueueFunc,
	}
}

func(h *configMapEventHandler)OnAdd(obj interface{}){
	cm,ok := obj.(*corev1.ConfigMap)
	if !ok {
		return
	}
	h.mayEnqueueZkWhenConfigMapChanged(cm)
}

func(h *configMapEventHandler)OnDelete(obj interface{}){
	cm,ok := obj.(*corev1.ConfigMap)
	if !ok {
		return
	}
	h.mayEnqueueZkWhenConfigMapChanged(cm)
}

func (h *configMapEventHandler)OnUpdate(oldObj, newObj interface{}){
	oldCm, ok := oldObj.(*corev1.ConfigMap)
	if !ok {
		return
	}
	newCm,ok := newObj.(*corev1.ConfigMap)
	if !ok {
		return
	}

	if newCm.ResourceVersion == oldCm.ResourceVersion{
		return
	}
	h.mayEnqueueZkWhenConfigMapChanged(newCm)
}


func(h *configMapEventHandler)mayEnqueueZkWhenConfigMapChanged(cm *corev1.ConfigMap){
	if len(cm.OwnerReferences) == 0{
		return
	}
	ownerReference := cm.OwnerReferences[0]
	zk,err := h.zkLister.ZookeeperClusters(cm.GetNamespace()).Get(ownerReference.Name)
	if err != nil{
		return
	}
	for _, ownerference := range cm.OwnerReferences{
		if ownerference.Kind != crdv1alpha1.ZookeeperKind || !metav1.IsControlledBy(cm, zk){
			return
		}
	}
	h.enqueueFunc(zk)
}

