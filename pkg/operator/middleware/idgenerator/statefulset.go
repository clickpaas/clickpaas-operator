package idgenerator

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
		ss := er.object.(*appv1.StatefulSet)
		return ss.DeepCopy(), nil
	case *crdv1alpha1.IdGenerate:
		redis := er.object.(*crdv1alpha1.IdGenerate)
		return newStatefulSetForIdGenerator(redis), nil
	}
	return nil, fmt.Errorf("trans object to service failed, unexcept type %#v", er.object)
}


func newStatefulSetForIdGenerator(generate *crdv1alpha1.IdGenerate)*appv1.StatefulSet{

	ss := &appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForIdGenerator(generate)},
			Namespace: generate.GetNamespace(),
			Name: getStatefulSetNameForIdGenerator(generate),
		},
		Spec: appv1.StatefulSetSpec{
			Replicas: &generate.Spec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: getLabelForIdGenerator(generate)},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: getLabelForIdGenerator(generate)},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: getStatefulSetNameForIdGenerator(generate),
							Image: generate.Spec.Image,
							ImagePullPolicy: corev1.PullPolicy(generate.Spec.ImagePullPolicy),
							Ports: []corev1.ContainerPort{
								{Name: "idg-port", ContainerPort: generate.Spec.Port},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: getVolumeDataName(generate.GetName()),
									MountPath: "/data/redis/16379",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: getVolumeDataName(generate.GetName()),
							VolumeSource: corev1.VolumeSource{EmptyDir: nil},
						},
					},
				},
			},
		},
	}
	return ss
}