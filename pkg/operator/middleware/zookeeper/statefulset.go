package zookeeper

import (
	"fmt"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

type statefulSetResourceEr struct {
	object interface{}
}

func(er *statefulSetResourceEr)StatefulSetResourceEr(...interface{})(*appv1.StatefulSet,error){
	switch er.object.(type) {
	case *appv1.StatefulSet:
		svc := er.object.(*appv1.StatefulSet)
		return svc.DeepCopy(), nil
	case *crdv1alpha1.ZookeeperCluster:
		zk := er.object.(*crdv1alpha1.ZookeeperCluster)
		return newStatefulSetForZookeeper(zk), nil
	}
	return nil, fmt.Errorf("unexcept type %#v", er.object)
}

func newStatefulSetForZookeeper(zk *crdv1alpha1.ZookeeperCluster)*appv1.StatefulSet{
	ss := &appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForZookeeperCluster(zk)},
			Name: getStatefulSetNameForZookeeper(zk),
			Namespace: zk.GetNamespace(),
		},
		Spec:       appv1.StatefulSetSpec{
			PodManagementPolicy: appv1.OrderedReadyPodManagement,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: getLabelForZookeeperCluster(zk),
				},
				Spec:       corev1.PodSpec{

				},
			},
		},
		Status:     appv1.StatefulSetStatus{},
	}
	return ss
}