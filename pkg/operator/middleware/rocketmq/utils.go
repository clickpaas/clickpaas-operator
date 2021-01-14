package rocketmq

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	"fmt"
)

func getStatefulSetNameForRocketmq(rocketmq *crdv1alpha1.Rocketmq)string{
	return fmt.Sprintf("%s", rocketmq.GetName())
}


func getDeploymentNameForRocketmq(rocketmq *crdv1alpha1.Rocketmq)string{
	return fmt.Sprintf("%s", rocketmq.GetName())
}


func getServiceNameForRocketmq(rocketmq *crdv1alpha1.Rocketmq)string{
	return fmt.Sprintf("%s", rocketmq.GetName())
}

func getServiceNameForRocketNameServer(rocketmq *crdv1alpha1.Rocketmq)string{
	return fmt.Sprintf("%s-nameserver", rocketmq.GetName())
}


func getLabelForRocketmqCluster(cluster *crdv1alpha1.Rocketmq)map[string]string{
	return map[string]string{"crdversion": crdv1alpha1.MiddlewareResourceVersion,
		"appname": cluster.GetName(),
		"kind": crdv1alpha1.RocketmqKind,
		"role": "rocketmq",
	}
}

func getLabelForRocketmqNameServer(rocketmq *crdv1alpha1.Rocketmq)map[string]string{
	return map[string]string{"crdversion": crdv1alpha1.MiddlewareResourceVersion,
		"appname": rocketmq.GetName(),
		"kind": crdv1alpha1.RocketmqKind,
		"role": "nameserver",
	}
}

func ownerReferenceForRocketmqCluster(cluster *crdv1alpha1.Rocketmq)metav1.OwnerReference{
	return *metav1.NewControllerRef(cluster, crdv1alpha1.SchemeGroupVersion.WithKind(crdv1alpha1.RocketmqKind))
}



func getVolumeNameForBrokerProperties(rocketmq *crdv1alpha1.Rocketmq)string{
	return fmt.Sprintf("%s-properties", rocketmq.GetName())
}

func getConfigMapNameForBrokerProperties(rocketmq *crdv1alpha1.Rocketmq)string{
	return fmt.Sprintf("%s", rocketmq.GetName())
}

func getBrokerPropertiesFileName(rocketmq *crdv1alpha1.Rocketmq)string{
	return fmt.Sprintf("%s.properties", rocketmq.GetName())
}