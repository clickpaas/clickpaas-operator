package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	RocketmqPlural = "rocketmqs"
	RocketmqSingular = "rocketmq"
	RocketmqShort = "mq"
	RocketmqKind = "Rocketmq"
)


// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Rocketmq represent the rocketmq application
type Rocketmq struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec RocketmqSpec `json:"spec"`
	Status RocketmqStatus `json:"status"`
}


// RocketmqSpec represent the spec of Rocketmq
type RocketmqSpec struct {
	CommonSpec `json:",inline"`
	ListenPort int32 `json:"listenPort"`
	HaPort int32 `json:"haPort"`
	FastPort int32 `json:"fastPort"`
	NameServerPort int32 `json:"nameServerPort"`

	Storage Storage `json:"storage"`
}


// RocketmqStatus represent the status of Rocketmq
type RocketmqStatus struct {

}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RocketmqList represent the list of Rocketmq
type RocketmqList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Rocketmq `json:"items"`
}