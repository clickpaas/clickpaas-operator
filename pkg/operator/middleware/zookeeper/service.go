package zookeeper

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

type serviceResourceEr struct {
	object interface{}
	f func(cluster *crdv1alpha1.ZookeeperCluster)*corev1.Service
}

func(er *serviceResourceEr)ServiceResourceEr(...interface{})(*corev1.Service,error){
	switch er.object.(type) {
	case *corev1.Service:
		cm := er.object.(*corev1.Service)
		return cm.DeepCopy(), nil
	case *crdv1alpha1.ZookeeperCluster:
		zkSvc := er.object.(*crdv1alpha1.ZookeeperCluster)
		return er.f(zkSvc), nil
	}
	return nil, fmt.Errorf("unknow type %#v", er.object)
}

func newServiceForZookeeperClient(cluster *crdv1alpha1.ZookeeperCluster)*corev1.Service{
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: getZookeeperClientServiceName(cluster),
			Namespace: cluster.GetNamespace(),
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForZookeeperCluster(cluster)},
		},
		Spec:       corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: getLabelForZookeeperCluster(cluster),
			Ports: []corev1.ServicePort{
				{Name: "zk-client", Protocol: corev1.ProtocolTCP, TargetPort: intstr.IntOrString{IntVal: cluster.Spec.ClientPort}, Port: cluster.Spec.ClientPort},
			},
		},
	}
	return svc
}


func newServiceForZookeeperServiceCommunicate(cluster *crdv1alpha1.ZookeeperCluster)*corev1.Service{
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: getZookeeperClusterCommunicateServiceName(cluster),
			Namespace: cluster.GetNamespace(),
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForZookeeperCluster(cluster)},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Selector: getLabelForZookeeperCluster(cluster),
			Ports: []corev1.ServicePort{
				{Name: "server-port", TargetPort: intstr.IntOrString{IntVal: cluster.Spec.ServerPort}, Port: cluster.Spec.ServerPort},
				{Name: "sync-port", TargetPort: intstr.IntOrString{IntVal: cluster.Spec.SyncPort}, Port: cluster.Spec.SyncPort},
			},
		},
	}
	return svc
}


