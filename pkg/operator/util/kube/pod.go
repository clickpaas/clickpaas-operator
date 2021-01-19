package kube

import (
	corev1 "k8s.io/api/core/v1"
	"time"
)

func FilteredActivePods(pods []*corev1.Pod)[]*corev1.Pod{

	// get all inactive pod from given a list of pods
	if pods == nil || len(pods) == 0{
		return nil
	}
	activePodList := make([]*corev1.Pod, 0)
	for _, podItem := range pods{
		// check status.phase of pod is running or not
		if podItem.Status.Phase != "Running"{continue}
		// check conditions, if condition's status is not True, ignore the pod
		conditionIsOk := true
		for _, condition := range podItem.Status.Conditions{
			if condition.Status != "True"{
				conditionIsOk = false
				break
			}
		}
		if ! conditionIsOk{continue}
		// check containerStatus
		for _,containerStatus := range podItem.Status.ContainerStatuses{
			if containerStatus.Ready != true{
				continue
			}
		}
		activePodList = append(activePodList, podItem)
	}
	return activePodList
}



func WaitForPodsReady(pods []*corev1.Pod, duration time.Duration)error{
	return nil
}