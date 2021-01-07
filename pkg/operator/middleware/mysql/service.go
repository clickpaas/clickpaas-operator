package mysql

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

func (er *serviceResourceEr)ServiceResourceEr(... interface{})(*corev1.Service,error){
	switch er.object.(type) {
	case *corev1.Service:
		svc := er.object.(*corev1.Service)
		return svc.DeepCopy(), nil
	case *crdv1alpha1.MysqlCluster:
		mysql := er.object.(*crdv1alpha1.MysqlCluster)
		return newServiceForMysql(mysql), nil
	}
	return nil, fmt.Errorf("unexcept type %#v", er.object)
}



func newServiceForMysql(cluster *crdv1alpha1.MysqlCluster)*corev1.Service{
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getServiceNameForMysql(cluster),
			Namespace: cluster.GetNamespace(),
			OwnerReferences: []metav1.OwnerReference{
				ownerReferenceForMysqlCluster(cluster),
			},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Ports: []corev1.ServicePort{
				{
					Name: "mysql-port",
					TargetPort: intstr.IntOrString{IntVal: cluster.Spec.Port},
					Port: cluster.Spec.Port,
				},
			},
			Selector: getLabelForMysqlCluster(cluster),
		},
	}
	return svc
}