package mongo

import (
	"fmt"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	"path"
	"strconv"
)

type statefulSetResourceEr struct {
	object interface{}
	nodeName string
}

func (er *statefulSetResourceEr)StatefulSetResourceEr(... interface{})(*appv1.StatefulSet,error){
	switch er.object.(type) {
	case *appv1.StatefulSet:
		ss := er.object.(*appv1.StatefulSet)
		return ss.DeepCopy(), nil
	case *crdv1alpha1.MongoCluster:
		mongo := er.object.(*crdv1alpha1.MongoCluster)
		return newStatefulSetForMongo(mongo, er.nodeName), nil
	}
	return nil, fmt.Errorf("unexcept type %#v", er.object)
}

func newStatefulSetForMongo(cluster *crdv1alpha1.MongoCluster, nodeName string)*appv1.StatefulSet{
	hostPathPolicy := corev1.HostPathDirectoryOrCreate
	ss := &appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:       cluster.GetNamespace(),
			Name:            getStatefulSetNameForMongoCluster(cluster),
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForMongoCluster(cluster)},
		},
		Spec: appv1.StatefulSetSpec{
			Replicas: &cluster.Spec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: getLabelForMongoCluster(cluster)},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: getLabelForMongoCluster(cluster),},
				Spec: corev1.PodSpec{
					NodeName: nodeName,
					Containers: []corev1.Container{
						{
							Name:            getStatefulSetNameForMongoCluster(cluster),
							ImagePullPolicy: corev1.PullPolicy(cluster.Spec.ImagePullPolicy),
							Image:           cluster.Spec.Image,
							Ports: []corev1.ContainerPort{
								{Name: "mongo-port", ContainerPort: cluster.Spec.Port},
							},
							Env: []corev1.EnvVar{
								{Name: "MONGO_INITDB_ROOT_USERNAME", Value: cluster.Spec.Config.User},
								{Name: "MONGO_INITDB_ROOT_PASSWORD", Value: cluster.Spec.Config.Password},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "vmongo-data",
									MountPath: "/data/db",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "vmongo-data",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Type: &hostPathPolicy,
									Path: path.Join("/data", cluster.GetName() , strconv.Itoa(int(cluster.Spec.Port))),
								},
							},
						},
					},
				},
			},
		},
	}
	return ss
}