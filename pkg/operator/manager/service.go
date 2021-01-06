package manager

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1lister "k8s.io/client-go/listers/core/v1"
	"l0calh0st.cn/clickpaas-operator/pkg/operator"
)

type serviceManager struct {
	operator.Manager
	serviceLister corev1lister.ServiceLister
}

func NewServiceManager(kubeClient kubernetes.Interface, serviceLister corev1lister.ServiceLister)*serviceManager{
	return &serviceManager{
		Manager:       operator.Manager{KubeClient: kubeClient},
		serviceLister: serviceLister,
	}
}


func (s *serviceManager) Create(obj interface{}, f func(interface{})(*corev1.Service,error)) (*corev1.Service, error) {
	svc,err := f(obj)
	if err != nil{
		return nil, err
	}
	return s.KubeClient.CoreV1().Services(svc.GetNamespace()).Create(context.TODO(), svc, metav1.CreateOptions{})
}

func (s *serviceManager) Update(obj interface{}, f func(interface{})(*corev1.Service,error)) (*corev1.Service, error) {
	svc,err := f(obj)
	if err != nil{
		return nil, err
	}
	return s.KubeClient.CoreV1().Services(svc.GetNamespace()).Update(context.TODO(), svc, metav1.UpdateOptions{})
}

func (s *serviceManager) Delete(obj interface{}, f func(interface{})(*corev1.Service,error)) error {
	svc,err := f(obj)
	if err != nil{
		return err
	}
	return s.KubeClient.CoreV1().Services(svc.GetNamespace()).Delete(context.TODO(), svc.GetName(), metav1.DeleteOptions{})
}

func (s *serviceManager) Get(obj interface{}, f func(interface{})(*corev1.Service,error)) (*corev1.Service, error) {
	svc,err := f(obj)
	if err != nil{
		return nil, err
	}
	return s.KubeClient.CoreV1().Services(svc.GetNamespace()).Get(context.TODO(), svc.GetName(), metav1.GetOptions{})
}


