package operator

import (
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

