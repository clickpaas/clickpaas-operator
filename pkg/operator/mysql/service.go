package mysql

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


func NewServiceManager(kubeClient kubernetes.Interface, svcLister corev1lister.ServiceLister)*serviceManager{
	return &serviceManager{kubeClient: kubeClient, serviceLister: svcLister}
}

func(m *serviceManager)Create(cluster *crdv1alpha1.MysqlCluster)(*corev1.Service,error){
	return m.kubeClient.CoreV1().Services(cluster.GetNamespace()).Create(context.TODO(), newServiceForMysql(cluster),metav1.CreateOptions{})
}


func(m *serviceManager)Get(cluster *crdv1alpha1.MysqlCluster)(*corev1.Service, error){
	return m.serviceLister.Services(cluster.GetNamespace()).Get(namedStatefulSetForMysql(cluster))
}


func(m *serviceManager)Update(svc *corev1.Service)(*corev1.Service, error){
	return m.kubeClient.CoreV1().Services(svc.GetNamespace()).Update(context.TODO(), svc, metav1.UpdateOptions{})
}

func(m *serviceManager)Delete(cluster *crdv1alpha1.MysqlCluster)error{
	return m.kubeClient.CoreV1().Services(cluster.GetNamespace()).Delete(context.TODO(), namedServiceForMysql(cluster), metav1.DeleteOptions{})
}

func newServiceForMysql(cluster *crdv1alpha1.MysqlCluster)*corev1.Service{
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: namedServiceForMysql(cluster),
			Namespace: cluster.GetNamespace(),
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "mysql-port", TargetPort: intstr.IntOrString{IntVal: cluster.Spec.Port}},
			},
			Selector: labelForMysqlCluster(cluster),
		},
	}
	return svc
}