package v1alpha1


import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	// CustomResourceName is the MySQLCluster custom resource definition qualified object name.
	RedisGCachePlural = "redisgcaches"
	RedisGCacheSingular = "redisgcache"
	RedisGCacheShort = "gcache"
	RedisGCacheKind = "RedisGCache"
)


// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RedisGCache represent the cluster of rediscluster
type RedisGCache struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec RedisGCacheSpec `json:"spec"`
	Status RedisGCacheStatus `json:"status"`
}


// RedisGCacheSpec represent the spec of RedisGCache
type RedisGCacheSpec struct {
	CommonSpec `json:",inline"`
	Port int32 `json:"port"`
	Storage Storage `json:"storage,omitempty"`
}


// RedisGCacheStatus represent the status of RedisGCache
type RedisGCacheStatus struct {

}


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RedisGCacheList represent the list of RedisGCache
type RedisGCacheList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []RedisGCache `json:"items"`
}



