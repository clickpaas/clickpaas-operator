package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	IdGeneratePlural = "idgenerates"
	IdGenerateSingular = "idgenerate"
	IdGenerateShort = "idgenerate"
	IdGenerateKind = "IdGenerate"
)


// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IdGenerate represent application idgenerate
type IdGenerate struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec IdGenerateSpec `json:"spec"`
	Status IdGenerateStatus `json:"status"`
}


// IdGenerateSpec represent the spec of IdGenerate
type IdGenerateSpec struct {
	CommonSpec `json:",inline"`
	Port int32 `json:"port"`
	Storage Storage `json:"storage"`
}

// IdGenerateStatus represent the status of IdGenerate
type IdGenerateStatus struct {

}


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IdGenerateList represent the list of IdGenerate
type IdGenerateList struct {
	metav1.TypeMeta	`json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	
	
	Items []IdGenerate `json:"items"`
}