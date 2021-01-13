package v1alpha1



type CommonSpec struct {
	Image string `json:"image,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
	Replicas int32 `json:"replicas,omitempty"`
}



type Storage struct {
	
}


type DatabaseViaConfig struct {
	User string `json:"user"`
	Password string `json:"password"`
	Host string `json:"host"`
	Port int `json:"port"`
}