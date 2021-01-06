package mysql

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

func getStatefulSetNameForMysql(cluster *crdv1alpha1.MysqlCluster)string{
	return fmt.Sprintf("%s-mysql", cluster.GetName())
}


func getConfigMapNameForMysql(cluster *crdv1alpha1.MysqlCluster)string{
	return fmt.Sprintf("%s-mysql", cluster.GetName())
}

func getServiceNameForMysql(cluster *crdv1alpha1.MysqlCluster)string{
	return fmt.Sprintf("%s-mysql", cluster.GetName())
}


func getLabelForMysqlCluster(cluster *crdv1alpha1.MysqlCluster)map[string]string{
	return map[string]string{"crdversion": crdv1alpha1.MiddlewareResourceVersion,
		"appname": cluster.GetName(),
		"kind": crdv1alpha1.MysqlKind,
	}
}


func ownerReferenceForMysqlCluster(cluster *crdv1alpha1.MysqlCluster)metav1.OwnerReference{
	return *metav1.NewControllerRef(cluster, crdv1alpha1.SchemeGroupVersion.WithKind(crdv1alpha1.MysqlKind))
}