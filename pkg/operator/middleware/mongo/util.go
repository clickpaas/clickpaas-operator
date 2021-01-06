package mongo

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	"fmt"
)

func getStatefulSetNameForMongoCluster(cluster *crdv1alpha1.MongoCluster)string{
	return fmt.Sprintf("%s-mongo", cluster.GetName())
}


func getConfigMapNameForMongoCluster(cluster *crdv1alpha1.MongoCluster)string{
	return fmt.Sprintf("%s-mongo", cluster.GetName())
}

func getServiceNameForMongo(cluster *crdv1alpha1.MongoCluster)string{
	return fmt.Sprintf("%s-mongo", cluster.GetName())
}


func getLabelForMongoCluster(cluster *crdv1alpha1.MongoCluster)map[string]string{
	return map[string]string{"crdversion": crdv1alpha1.MiddlewareResourceVersion,
		"appname": cluster.GetName(),
		"kind": crdv1alpha1.MongoClusterKind,
	}
}


func ownerReferenceForMongoCluster(cluster *crdv1alpha1.MongoCluster)metav1.OwnerReference{
	return *metav1.NewControllerRef(cluster, crdv1alpha1.SchemeGroupVersion.WithKind(crdv1alpha1.MongoClusterKind))
}
