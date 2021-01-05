package mysql

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	"l0calh0st.cn/clickpaas-operator/pkg/crd/middleware"
)

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
	return map[string]string{"crdversion": crdv1alpha1.MiddlewareResourceVersion, "appname": cluster.GetName()}
}


func ownerReferenceForMysqlCluster(cluster *crdv1alpha1.MysqlCluster)metav1.OwnerReference{
	return *metav1.NewControllerRef(cluster, crdv1alpha1.SchemeGroupVersion.WithKind(middleware.MysqlKind))
}