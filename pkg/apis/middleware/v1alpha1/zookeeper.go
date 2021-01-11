package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"



const (
	ZookeeperPlural = "zookeeperclusters"
	ZookeeperSingular = "zookeepercluster"
	ZookeeperShort = "zk"
	ZookeeperKind = "ZookeeperCluster"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ZookeeperCluster struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	
	Spec ZookeeperClusterSpec `json:"spec"`
	Status ZookeeperClusterStatus `json:"status"`
}


// ZookeeperClusterSpec represent the Spec of ZookeeperCluster
type ZookeeperClusterSpec struct {
	CommonSpec `json:",inline"`	
	ClientPort int32 `json:"clientPort"`
	SyncPort int32 `json:"syncPort,omitempty"`
	ServerPort int32 `json:"ServerPort,omitempty"`
	Storage Storage `json:"storage,omitempty"`
	Command []string `json:"args,omitempty"`
	Args []string `json:"args,omitempty"`
}


// ZookeeperClusterConfig represent the config of ZookeeperCluster
type ZookeeperClusterConfig struct {
	
}

// ZookeeperClusterStatus represent the Status of ZookeeperCluster
type ZookeeperClusterStatus struct {
	
}




// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ZookeeperClusterList represent the list of ZookeeperCluster
type ZookeeperClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ZookeeperCluster `json:"items"`
}
