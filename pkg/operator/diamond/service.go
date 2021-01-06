package diamond

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	corev1lister "k8s.io/client-go/listers/core/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

type serviceManager struct {
	kubeClient kubernetes.Interface
	serviceLister corev1lister.ServiceLister
}

func NewServiceManager(kubeClient kubernetes.Interface, serviceLister corev1lister.ServiceLister)*serviceManager{
	return &serviceManager{
		kubeClient:    kubeClient,
		serviceLister: serviceLister,
	}
}


func(m *serviceManager)Get(diamond *crdv1alpha1.Diamond)(*corev1.Service,error){
	return m.serviceLister.Services(diamond.GetNamespace()).Get(getServiceNameForDiamond(diamond))
}


func(m *serviceManager)Delete(diamond *crdv1alpha1.Diamond)error{
	return m.kubeClient.CoreV1().Services(diamond.GetNamespace()).Delete(context.TODO(), getServiceNameForDiamond(diamond), metav1.DeleteOptions{})
}

func(m *serviceManager)Update(svc *corev1.Service)(*corev1.Service, error){
	return m.kubeClient.CoreV1().Services(svc.GetNamespace()).Update(context.TODO(), svc, metav1.UpdateOptions{})
}

func(m *serviceManager)Create(diamond *crdv1alpha1.Diamond)(*corev1.Service, error){
	return m.kubeClient.CoreV1().Services(diamond.GetNamespace()).Create(context.TODO(), newServiceForDiamond(diamond), metav1.CreateOptions{})
}

func newServiceForDiamond(diamond *crdv1alpha1.Diamond)*corev1.Service{
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: getServiceNameForDiamond(diamond),
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForDiamond(diamond)},
			Namespace: diamond.GetNamespace(),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: getLabelForDiamond(diamond),
			Ports: []corev1.ServicePort{
				{Name: "diamond-port", Port: diamond.Spec.Port, TargetPort: intstr.IntOrString{IntVal: diamond.Spec.Port}},
			},
		},
	}
	return svc
}