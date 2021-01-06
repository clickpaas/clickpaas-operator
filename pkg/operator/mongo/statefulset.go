package mongo

import (
	"fmt"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)


func statefulSetObjHandleFunc(obj interface{})(*appv1.StatefulSet,error){
	switch obj.(type) {
	case *appv1.StatefulSet:
		ss := obj.(*appv1.StatefulSet)
		return ss.DeepCopy(), nil
	case *crdv1alpha1.MongoCluster:
		mongo := obj.(*crdv1alpha1.MongoCluster)
		return newStatefulSetForMongo(mongo), nil
	}
	return nil, fmt.Errorf("unexcept type %#v", obj)
}

func newStatefulSetForMongo(cluster *crdv1alpha1.MongoCluster)*appv1.StatefulSet{
	ss := &appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.GetNamespace(),
			Name: getStatefulSetNameForMongoCluster(cluster),
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForMongoCluster(cluster)},
		},
		Spec: appv1.StatefulSetSpec{
			Replicas: &cluster.Spec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: getLabelForMongoCluster(cluster)},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: getLabelForMongoCluster(cluster),},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: getStatefulSetNameForMongoCluster(cluster),
							ImagePullPolicy: corev1.PullPolicy(cluster.Spec.ImagePullPolicy),
							Image: cluster.Spec.Image,
							Ports: []corev1.ContainerPort{
								{Name: "mongo-port", ContainerPort: cluster.Spec.Port},
							},
							Env: []corev1.EnvVar{
								{Name: "MONGO_INITDB_ROOT_USERNAME", Value: cluster.Spec.Config.User},
								{Name: "MONGO_INITDB_ROOT_PASSWORD", Value: cluster.Spec.Config.Password},
							},
						},
					},
				},
			},
		},
	}
	return ss
}