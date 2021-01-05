package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:defaulter-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MongoCluster struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	
	Spec MongoClusterSpec `json:"spec"`
	Status MongoClusterStatus `json:"status"`
	
}


// MongoClusterSpec represent the spec of MongoCluster
type MongoClusterSpec struct {
	CommonSpec `json:",inline"`
	Port int32 `json:"port"`
	Storage Storage `json:"storage"`
	Config MongoClusterConfig `json:"config"`
}

// MongoClusterStatus represent the status of MongoCluster
type MongoClusterStatus struct {
	
}

type MongoClusterConfig struct {
	User string `json:"user"`
	Password string `json:"password"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MongoClusterList represent the list of MongoCluster
type MongoClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	
	Items []MongoCluster
}