package diamond

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
func (s *serviceResourceEr)ServiceResourceEr(...interface{})(*corev1.Service,error){
	switch s.object.(type) {
	case *corev1.Service:
		svc := s.object.(*corev1.Service)
		return svc.DeepCopy(), nil
	case *crdv1alpha1.Diamond:
		diamond := s.object.(*crdv1alpha1.Diamond)
		return newServiceForDiamond(diamond), nil
	}
	return nil, fmt.Errorf("unexcept type %#v", s.object)
}


func newServiceForDiamond(diamond *crdv1alpha1.Diamond)*corev1.Service{
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            getServiceNameForDiamond(diamond),
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForDiamond(diamond)},
			Namespace:       diamond.GetNamespace(),
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: getLabelForDiamond(diamond),
			Ports: []corev1.ServicePort{
				{Name: "diamond-port", Port: diamond.Spec.Port, TargetPort: intstr.IntOrString{IntVal: diamond.Spec.Port}},
			},
		},
	}
	return svc
}