package gcache

import (
	"fmt"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

func statefulSetResourceHandleFunc(obj interface{})(*appv1.StatefulSet, error){
	switch obj.(type) {
	case *appv1.StatefulSet:
		ss := obj.(*appv1.StatefulSet)
		return ss.DeepCopy(), nil
	case *crdv1alpha1.RedisGCache:
		gcache := obj.(*crdv1alpha1.RedisGCache)
		return newStatefulSetForRedisGCache(gcache), nil
	}
	return nil, fmt.Errorf("trans object to statefulset failed, unexcept type %#v", obj)
}

func newStatefulSetForRedisGCache(redis *crdv1alpha1.RedisGCache)*appv1.StatefulSet{
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
						},
					},
				},
			},
		},
	}
	return ss
}
