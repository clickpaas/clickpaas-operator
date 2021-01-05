package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// service configmap statefulset

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Diamond represent diamond cluster
type Diamond struct {
	metav1.TypeMeta `json:",inline"`	
	metav1.ObjectMeta `json:"metadata,omitempty"`
	
	Spec DiamondSpec `json:"spec"`
	Status DiamondStatus `json:"status"`
}

// DiamondSpec represent the spec of diamond
type DiamondSpec struct {
	CommonSpec `json:",inline"`
	Port int32 `json:"port"`
	Config DatabaseViaConfig `json:"config"`
	Storage Storage `json:"storage,omitempty"`
}


// DiamondStatus represent the status of diamond
type DiamondStatus struct {
	
}


// DiamondConfig represent the config of diamond
type DiamondConfig struct {
	
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DiamondList represent the list of Diamond
type DiamondList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	
	Items []Diamond `json:"items"`
}


