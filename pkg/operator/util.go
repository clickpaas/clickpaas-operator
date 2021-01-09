package operator

import (
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func WaitForStatefulSetPodsReady(ss *appv1.StatefulSet,timeout time.Duration)error{
	return nil
}


func WaitForDeploymentPodsReady(deploy *appv1.Deployment,timeout time.Duration)error{
	return nil
}

func WaitForSinglePodReady(pod *corev1.Pod, timeout time.Duration)error{
	return nil
}


func WaitAllPodsReady(podList []*corev1.Pod)bool{
	return true
}


func AddOwnerReference(obj metav1.Object){

}



func PodIsReadyOrNot(pod *corev1.Pod)bool{
	if pod == nil{
		return false
	}
	for _,condition :=range  pod.Status.Conditions{
		if condition.Status != "True"{
			return false
		}
	}
	for _, container := range pod.Status.ContainerStatuses{
		if container.Ready != true{
			return false
		}
	}
	return true
}