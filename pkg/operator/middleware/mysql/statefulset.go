package mysql

import (
	"fmt"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

type statefulSetEr struct {
	object interface{}
}

func (er *statefulSetEr)StatefulSetResourceEr(... interface{})(*appv1.StatefulSet,error){
	switch er.object.(type) {
	case *appv1.StatefulSet:
		ss := er.object.(*appv1.StatefulSet)
		return ss.DeepCopy(), nil
	case *crdv1alpha1.MysqlCluster:
		mysql := er.object.(*crdv1alpha1.MysqlCluster)
		return newStatefulSetForMysqlCluster(mysql), nil
	}
	return nil, fmt.Errorf("unexcept type %#v", er.object)
}


func newStatefulSetForMysqlCluster(cluster *crdv1alpha1.MysqlCluster)*appv1.StatefulSet{
	var hostPathCreateIfNotExisted = corev1.HostPathDirectoryOrCreate
	ss := &appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:            getStatefulSetNameForMysql(cluster),
			Namespace:       cluster.GetNamespace(),
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForMysqlCluster(cluster)},
		},
		Spec: appv1.StatefulSetSpec{
			Replicas: &cluster.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: getLabelForMysqlCluster(cluster),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: getLabelForMysqlCluster(cluster)},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyAlways,
					Containers: []corev1.Container{
						{
							ImagePullPolicy: corev1.PullPolicy(cluster.Spec.ImagePullPolicy),
							Env: []corev1.EnvVar{
								{Name: "MYSQL_ROOT_PASSWORD", Value: cluster.Spec.Config.Password},
								{Name: "MYSQL_ROOT_HOST", Value: "%"},
							},
							Args: cluster.Spec.Args,
							Ports: []corev1.ContainerPort{
								{Name: "mysql-port", ContainerPort: cluster.Spec.Port},
							},
							Image: cluster.Spec.Image,
							Name:  getStatefulSetNameForMysql(cluster),
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "bootstrap",
									ReadOnly: true,
									MountPath: "/tmp/lib",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "bootstrap",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/data/lib/",
									Type: &hostPathCreateIfNotExisted,
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