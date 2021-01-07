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
	
	Spec LtsJobTrackerSpec `json:"status"`
	Status LtsJobTrackerStatus `json:"status"`
}

// LtsJobTrackerSpec represent the spec of LtsJobTracker
type LtsJobTrackerSpec struct {
	CommonSpec `json:",inline"`
	Port int32 `json:"port"`
	HealthPort int32 `json:"healthPort"`
	Storage Storage `json:"storage"`
	Config LtsJobTrackerConfig `json:"config"`
}


type LtsJobTrackerConfig struct {
	Db DatabaseViaConfig `json:"db"`
	RegistryAddress string `json:"registryAddress"`
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



