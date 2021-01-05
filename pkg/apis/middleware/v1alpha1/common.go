package v1alpha1



type CommonSpec struct {
	Image string `json:"image"`
	ImagePullPolicy string `json:"imagePullPolicy"`
	Replicas int32 `json:"replicas"`
}



type Storage struct {
	
}


type DatabaseViaConfig struct {
	User string `json:"user"`
	Password string `json:"password"`
	Host string `json:"host"`
	Port string `json:"port"`
}