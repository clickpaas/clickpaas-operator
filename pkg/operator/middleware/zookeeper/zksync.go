package zookeeper

import (
	corev1 "k8s.io/api/core/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	"strconv"
)


type zookeeperPodList []zookeeperPod


type zookeeperPod struct {
	name string
	id string
}

func(zp zookeeperPod)isEqual(other zookeeperPod)bool{
	if zp.name == other.name && zp.id == other.id{
		return true
	}
	return false
}

func getShouldBeDeletedPodList(kubePods []*corev1.Pod, zkPods []zookeeperPod)[]*corev1.Pod{
	var unExcepted []*corev1.Pod
	if len(zkPods) == 0{
		// 如果replication个数是0，则所有相关的pod 都应该被删除
		for _, pod := range kubePods{
			unExcepted = append(unExcepted, pod)
		}
		return unExcepted
	}
	for _, pod := range kubePods{
		for _,zkPod := range zkPods{
			if kubePodConversionToZkPod(pod).isEqual(zkPod){
				goto CONTINUE
			}
		}
		unExcepted = append(unExcepted, pod)
	CONTINUE:
	}
	return unExcepted
}

func getShouldInstalledPodList(kubePodS []*corev1.Pod, zkPods []zookeeperPod)[]zookeeperPod{
	var shouldInstall []zookeeperPod
	for _, zkPod := range zkPods{
		for _, pod := range kubePodS{
			if kubePodConversionToZkPod(pod).isEqual(zkPod){
				goto CONTINUE
			}
		}
		shouldInstall = append(shouldInstall, zkPod)
	CONTINUE:
	}
	return shouldInstall
}


func kubePodConversionToZkPod(pod *corev1.Pod)zookeeperPod{
	return zookeeperPod{
		name: pod.GetName(),
		id:   pod.GetAnnotations()["id"],
	}
}


func generateZookeeperPodList(cluster *crdv1alpha1.ZookeeperCluster)zookeeperPodList{
	zkPodList := []zookeeperPod{}
	for i := 0; i < int(cluster.Spec.Replicas); i++{
		zkPodList = append(zkPodList, zookeeperPod{name: getSinglePodNameForZookeeperCluster(cluster, i), id: strconv.Itoa(i)})
	}
	return zkPodList
}