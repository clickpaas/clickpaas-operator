package idgenerator

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

type serviceResourceEr struct {
	object interface{}
}

func(er *serviceResourceEr)ServiceResourceEr(...interface{})(*corev1.Service,error){
	switch er.object.(type) {
	case *corev1.Service:
		svc := er.object.(*corev1.Service)
		return svc.DeepCopy(), nil
	case *crdv1alpha1.IdGenerate:
		redis := er.object.(*crdv1alpha1.IdGenerate)
		return newServiceForIdGenerator(redis), nil
	}
	return nil, fmt.Errorf("trans object to service failed, unexcept type %#v", er.object)
}


func newServiceForIdGenerator(generate *crdv1alpha1.IdGenerate)*corev1.Service{
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            getServiceNameForIdGenerator(generate),
			Namespace:       generate.GetNamespace(),
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForIdGenerator(generate)},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Selector: getLabelForIdGenerator(generate),
			Ports: []corev1.ServicePort{
				{Name: "idgenerator-port", TargetPort: intstr.IntOrString{IntVal: generate.Spec.Port}, Port: generate.Spec.Port},
			},
		},
	}
	return svc
}