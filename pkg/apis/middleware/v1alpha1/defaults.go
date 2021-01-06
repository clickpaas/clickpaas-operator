package v1alpha1



// WithDefaultsMysqlCluster default not specified fields
func WithDefaultsMysqlCluster(cluster *MysqlCluster){
	if len(cluster.Spec.Args) == 0{
		cluster.Spec.Args = []string{"--character-set-server=utf8mb4", "--collation-server=utf8mb4_unicode_ci", "--lower_case_table_names=1"}
	}
	if cluster.Spec.Config.Password == "" {
		cluster.Spec.Config.Password = "diamond^^^"
	}
	if cluster.Spec.Config.User == ""{
		cluster.Spec.Config.User = "root"
	}
	if cluster.Spec.ImagePullPolicy == ""{
		cluster.Spec.ImagePullPolicy = "IfNotPresent"
	}
	if cluster.Spec.Replicas == 0{
		cluster.Spec.Replicas = 1
	}
	if cluster.Spec.Image == ""{
		cluster.Spec.Image = "registry.bizsaas.net/mysql:5.7.22"
	}
	if cluster.Spec.Port == 0{
		cluster.Spec.Port = 3306
	}
	// no
}


// WithDefaultsDiamond set default value for not special fields
func WithDefaultsDiamond(diamond *Diamond){
	if diamond.Spec.Port == 0{
		diamond.Spec.Port = 80
	}
	if diamond.Spec.Config.User == ""{
		diamond.Spec.Config.User = "root"
	}
	if diamond.Spec.Config.Password == ""{
		diamond.Spec.Config.Password = "diamond^^^"
	}
	if diamond.Spec.Replicas == 0{
		diamond.Spec.Replicas = 1
	}
	if diamond.Spec.ImagePullPolicy == ""{
		diamond.Spec.ImagePullPolicy = "IfNotPresent"
	}
	if diamond.Spec.Image == ""{
		diamond.Spec.Image = "registry.bizsaas.net/diamond:2.0.0-r2"
	}
}