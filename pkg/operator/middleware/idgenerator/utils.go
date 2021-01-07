package idgenerator

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	"fmt"
)

func getStatefulSetNameForIdGenerator(redis *crdv1alpha1.IdGenerate)string{
	return fmt.Sprintf("%s-diamond", redis.GetName())
}

func getConfigMapNameForIdGenerator(redis *crdv1alpha1.IdGenerate)string{
	return fmt.Sprintf("%s-diamond", redis.GetName())
}

func getServiceNameForIdGenerator(redis *crdv1alpha1.IdGenerate)string{
	return fmt.Sprintf("%s-diamond", redis.GetName())
}


func getLabelForIdGenerator(redis *crdv1alpha1.IdGenerate)map[string]string{
	return map[string]string{"crdversion": crdv1alpha1.MiddlewareResourceVersion,
		"appname": redis.GetName(),
		"kind": crdv1alpha1.IdGenerateKind,
	}
}

func ownerReferenceForIdGenerator(redis *crdv1alpha1.IdGenerate)metav1.OwnerReference{
	return *metav1.NewControllerRef(redis, crdv1alpha1.SchemeGroupVersion.WithKind(crdv1alpha1.IdGenerateKind))
}
