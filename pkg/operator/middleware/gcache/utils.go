package gcache

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	"fmt"
)

func getStatefulSetNameForRedisGCache(redis *crdv1alpha1.RedisGCache)string{
	return fmt.Sprintf("%s", redis.GetName())
}

func getConfigMapNameForRedisGCache(redis *crdv1alpha1.RedisGCache)string{
	return fmt.Sprintf("%s", redis.GetName())
}

func getServiceNameForRedisGCache(redis *crdv1alpha1.RedisGCache)string{
	return fmt.Sprintf("%s", redis.GetName())
}


func getLabelForRedisGCache(redis *crdv1alpha1.RedisGCache)map[string]string{
	return map[string]string{"crdversion": crdv1alpha1.MiddlewareResourceVersion,
		"appname": redis.GetName(),
		"kind": crdv1alpha1.RedisGCacheKind,
	}
}

func ownerReferenceForRedisGCache(redis *crdv1alpha1.RedisGCache)metav1.OwnerReference{
	return *metav1.NewControllerRef(redis, crdv1alpha1.SchemeGroupVersion.WithKind(crdv1alpha1.RedisGCacheKind))
}


func getMountPathForData(podName string)string{
	return fmt.Sprintf("%s-data", podName)
}
