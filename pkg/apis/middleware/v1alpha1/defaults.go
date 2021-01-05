package v1alpha1


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


