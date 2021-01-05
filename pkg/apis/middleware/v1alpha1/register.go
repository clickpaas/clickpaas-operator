package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"l0calh0st.cn/clickpaas-operator/pkg/apis/middleware"
)

const (
	MiddlewareResourceVersion = "v1alpha1"
)

var (
	SchemeGroupVersion = schema.GroupVersion{
		Group:   middleware.GroupName,
		Version: MiddlewareResourceVersion,
	}
)

var (
	SchemeBuilder  runtime.SchemeBuilder
	localSchemeBuilder = &SchemeBuilder
	AddToScheme = localSchemeBuilder.AddToScheme
)

func Resource(resource string)schema.GroupResource{
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}


func Kind(kind string)schema.GroupKind{
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

func init(){
	localSchemeBuilder.Register(addKnownTypes)
}

func addKnownTypes(scheme *runtime.Scheme)error{
	scheme.AddKnownTypes(SchemeGroupVersion,
		new(Diamond),
		new(MysqlCluster),
		new(MongoCluster),
		new(ZookeeperCluster),
		new(Rocketmq),
		new(LtsJobTracker),
		new(RedisGCache),
		new(IdGenerate))
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
