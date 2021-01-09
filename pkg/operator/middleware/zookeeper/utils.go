package zookeeper

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	"strconv"
)

func getStatefulSetNameForZookeeper(cluster *crdv1alpha1.ZookeeperCluster)string{
	return fmt.Sprintf("%s-mysql", cluster.GetName())
}


func getConfigMapNameForZookeeper(cluster *crdv1alpha1.ZookeeperCluster)string{
	return fmt.Sprintf("%s-mysql", cluster.GetName())
}


func getZookeeperClientServiceName(cluster *crdv1alpha1.ZookeeperCluster)string{
	return fmt.Sprintf("%s-zk-client", cluster.GetName())
}

func getZookeeperClusterCommunicateServiceName(cluster *crdv1alpha1.ZookeeperCluster)string{
	return fmt.Sprintf("%s-zk-server", cluster.GetName())
}

func getLabelForZookeeperCluster(cluster *crdv1alpha1.ZookeeperCluster)map[string]string{
	return map[string]string{"crdversion": crdv1alpha1.MiddlewareResourceVersion,
	"appname": cluster.GetName(),
	"kind": crdv1alpha1.ZookeeperKind,
	}
}


func ownerReferenceForZookeeperCluster(cluster *crdv1alpha1.ZookeeperCluster)metav1.OwnerReference{
	return *metav1.NewControllerRef(cluster, crdv1alpha1.SchemeGroupVersion.WithKind(crdv1alpha1.ZookeeperKind))
}

func getSinglePodNameForZookeeperCluster(cluster *crdv1alpha1.ZookeeperCluster, id int)string{
	return fmt.Sprintf("%s-%d", cluster.GetName(), id)
}

func getPodSubDomainForZookeeperCluster(cluster *crdv1alpha1.ZookeeperCluster)string{
	return fmt.Sprintf("%s-subdomain", cluster.GetName())
}

func getPodHostNameForZookeeperCluster(cluster *crdv1alpha1.ZookeeperCluster, id int)string{
	return fmt.Sprintf("%s-%d", cluster.GetName(), id)
}


func getOrGenerateZkServerPodSName(cluster *crdv1alpha1.ZookeeperCluster)[]string{
	var nameList []string
	for i := 0; i < int(cluster.Spec.Replicas); i++{
		nameList = append(nameList, getSinglePodNameForZookeeperCluster(cluster, i))
	}
	return nameList
}

func getAnnotationsForPod(id int)map[string]string{
	return map[string]string{"id": strconv.Itoa(id)}
}