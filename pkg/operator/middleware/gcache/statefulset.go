package gcache

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

func (ss *statefulSetResourceEr)StatefulSetResourceEr(... interface{})(*appv1.StatefulSet, error){
	switch ss.object.(type) {
	case *appv1.StatefulSet:
		ss := ss.object.(*appv1.StatefulSet)
		return ss.DeepCopy(), nil
	case *crdv1alpha1.RedisGCache:
		gcache := ss.object.(*crdv1alpha1.RedisGCache)
		return newStatefulSetForRedisGCache(gcache, ss.nodeName), nil
	}
	return nil, fmt.Errorf("trans object to statefulset failed, unexcept type %#v", ss.object)
}

func newStatefulSetForRedisGCache(redis *crdv1alpha1.RedisGCache, nodeName string)*appv1.StatefulSet{
	hostPathPolicy := corev1.HostPathDirectoryOrCreate
	ss := &appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForRedisGCache(redis)},
			Namespace: redis.GetNamespace(),
			Name: getStatefulSetNameForRedisGCache(redis),
		},
		Spec: appv1.StatefulSetSpec{
			Replicas: &redis.Spec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: getLabelForRedisGCache(redis)},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: getLabelForRedisGCache(redis)},
				Spec: corev1.PodSpec{
					NodeName: nodeName,
					Containers: []corev1.Container{
						{
							Name: getStatefulSetNameForRedisGCache(redis),
							Env: []corev1.EnvVar{
								{Name: "NETWORK_INTERFACE", Value: "eth0"},
							},
							Image: redis.Spec.Image,
							ImagePullPolicy: corev1.PullPolicy(redis.Spec.ImagePullPolicy),
							Ports: []corev1.ContainerPort{
								{Name: "redis-port", ContainerPort: redis.Spec.Port},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: getMountPathForData(redis.GetName()),
									MountPath: "/data/redis/8300",
								},
								{
									Name: "vredis",
									MountPath: "/data/redis/8300",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: getMountPathForData(redis.GetName()),
							VolumeSource: corev1.VolumeSource{EmptyDir: nil},
						},
						{
							Name: "vredis",
							VolumeSource:corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: path.Join("data", redis.GetName(), strconv.Itoa(int(redis.Spec.Port))),
									Type: &hostPathPolicy,
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
