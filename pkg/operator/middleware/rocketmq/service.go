package rocketmq

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

type serviceResourceEr struct {
	object interface{}
	f func(rocketmq *crdv1alpha1.Rocketmq)*corev1.Service
}

func(er *serviceResourceEr)ServiceResourceEr(...interface{})(*corev1.Service,error){
	switch er.object.(type) {
	case *corev1.Service:
		svc := er.object.(*corev1.Service)
		return svc.DeepCopy(), nil
	case *crdv1alpha1.Rocketmq:
		rocketmq := er.object.(*crdv1alpha1.Rocketmq)
		return er.f(rocketmq), nil
	}
	return nil, fmt.Errorf("unexcept type %#v", er.object)
}

func newServiceForRocketmqNameServer(rocketmq *crdv1alpha1.Rocketmq)*corev1.Service{
	nsSvc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForRocketmqCluster(rocketmq)},
			Name: getServiceNameForRocketNameServer(rocketmq),
			Namespace: rocketmq.GetNamespace(),
		},
		Spec:       corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: getLabelForRocketmqNameServer(rocketmq),
			Ports: []corev1.ServicePort{
				{Name: "ns-port", TargetPort: intstr.IntOrString{IntVal: rocketmq.Spec.NameServerPort}, Port: rocketmq.Spec.NameServerPort},
			},
		},
	}
	return nsSvc
}




func newServiceForRocketmq(rocketmq *crdv1alpha1.Rocketmq)*corev1.Service{
	rkSvc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForRocketmqCluster(rocketmq)},
			Name: getServiceNameForRocketmq(rocketmq),
			Namespace: rocketmq.GetNamespace(),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Selector: getLabelForRocketmqCluster(rocketmq),
			Ports: []corev1.ServicePort{
				{Name: "fport", TargetPort: intstr.IntOrString{IntVal: rocketmq.Spec.FastPort}, Port: rocketmq.Spec.FastPort},
				{Name: "haport", TargetPort: intstr.IntOrString{IntVal: rocketmq.Spec.HaPort}, Port: rocketmq.Spec.HaPort},
				{Name:"client-port", TargetPort: intstr.IntOrString{IntVal: rocketmq.Spec.ListenPort}, Port: rocketmq.Spec.ListenPort},
			},
		},
	}
	return rkSvc
}