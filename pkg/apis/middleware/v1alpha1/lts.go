package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	LtsJobTrackerPlural = "ltsjobtrackers"
	LtsJobTrackerSingular = "ltsjobtracker"
	LtsJobTrackerShort = "lts"
	LtsJobTrackerKind = "LtsJobTracker"
)


// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LtsJobTracker represent the lts application
type LtsJobTracker struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	
	Spec LtsJobTrackerSpec `json:"spec"`
	Status LtsJobTrackerStatus `json:"status"`
}

// LtsJobTrackerSpec represent the spec of LtsJobTracker
type LtsJobTrackerSpec struct {
	CommonSpec `json:",inline"`
	Port int32 `json:"port,omitempty"`
	HealthPort int32 `json:"healthPort,omitempty"`
	Storage Storage `json:"storage,omitempty"`
	Config LtsJobTrackerConfig `json:"config,omitempty"`
}


type LtsJobTrackerConfig struct {
	Db DatabaseViaConfig `json:"db,omitempty"`
	RegistryAddress string `json:"registryAddress,omitempty"`
}

// LtsJobTrackerStatus represent the status of LtsJobTracker
type LtsJobTrackerStatus struct {

}


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LtsJobTrackerList represent the list of LtsJobTracker
type LtsJobTrackerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	
	Items []LtsJobTracker `json:"items"`
}



