package mysql

import crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"

func namedStatefulSetForMysql(cluster *crdv1alpha1.MysqlCluster)string{
	return cluster.GetName()
}


func namedConfigMapForMysql(cluster *crdv1alpha1.MysqlCluster)string{
	return cluster.GetName()
}

func namedServiceForMysql(cluster *crdv1alpha1.MysqlCluster)string{
	return cluster.GetName()
}


func labelForMysqlCluster(cluster *crdv1alpha1.MysqlCluster)map[string]string{
	return map[string]string{}
}