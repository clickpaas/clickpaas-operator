package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	MysqlClusterPlural = "mysqlclusters"
	MysqlSingular = "mysqlcluster"
	MysqlShort = "mysql"
	MysqlKind = "MysqlCluster"
)


// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MysqlCluster represent the mysql cluster
type MysqlCluster struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec MysqlClusterSpec `json:"spec"`
	Status MysqlClusterStatus `json:"status"`

}

// MysqlCluster
type MysqlClusterSpec struct {
	CommonSpec `json:",inline"`
	Port int32 `json:"port"`
	Args []string `json:"args,omitempty"`
	Config MysqlClusterConfig `json:"config"`
	Storage Storage `json:"storage,omitempty"`
}

// MySqlClusterStatus represent the status of MySqlCluster
type MysqlClusterStatus struct {

}


// MysqlClusterConfig represent the config of MySqlCluster
type MysqlClusterConfig struct {
	User string `json:"user"`
	Password string `json:"password"`
}


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MysqlClusterList represent the list of MysqlCluster
type MysqlClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []MysqlCluster `json:"items"`
}