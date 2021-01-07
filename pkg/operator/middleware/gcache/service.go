package gcache

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
	case *crdv1alpha1.RedisGCache:
		redis := er.object.(*crdv1alpha1.RedisGCache)
		return newServiceForRedisGCache(redis), nil
	}
	return nil, fmt.Errorf("trans object to service failed, unexcept type %#v", er.object)
}

func newServiceForRedisGCache(redisGCache *crdv1alpha1.RedisGCache)*corev1.Service{
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            getServiceNameForRedisGCache(redisGCache),
			Namespace:       redisGCache.GetNamespace(),
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForRedisGCache(redisGCache)},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Selector: getLabelForRedisGCache(redisGCache),
			Ports: []corev1.ServicePort{
				{Name: "redis-port", TargetPort: intstr.IntOrString{IntVal: redisGCache.Spec.Port}, Port: redisGCache.Spec.Port},
			},
		},
	}
	return svc
}
