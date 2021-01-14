package rocketmq

import (
	"fmt"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	"l0calh0st.cn/clickpaas-operator/pkg/operator/middleware/rocketmq/res"
)

type configMapEr struct {
	object interface{}
}

func(er *configMapEr)ConfigMapResourceEr(...interface{})(*corev1.ConfigMap,error){
	switch er.object.(type) {
	case *appv1.Deployment:
		svc := er.object.(*corev1.ConfigMap)
		return svc.DeepCopy(), nil
	case *crdv1alpha1.Rocketmq:
		rocketmq := er.object.(*crdv1alpha1.Rocketmq)
		return newConfigMapForSampleBrokerProperties(rocketmq), nil
	}
	return nil, fmt.Errorf("unexcept type %#v", er.object)
}


func newConfigMapForSampleBrokerProperties(rocketmq *crdv1alpha1.Rocketmq)*corev1.ConfigMap{
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: getConfigMapNameForBrokerProperties(rocketmq),
			Namespace: rocketmq.GetNamespace(),
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForRocketmqCluster(rocketmq)},
		},
		Data: map[string]string{getBrokerPropertiesFileName(rocketmq): res.NewSampleBrokerProperties(rocketmq, res.BrokerRoleSyncMaster)},
		BinaryData: nil,
	}
	return cm
}