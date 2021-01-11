package zookeeper

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	"strconv"
)




type podResourceEr struct {
	object interface{}
	id int
}
func deepCopyPodResource(pod *corev1.Pod)*podResourceEr{
	return &podResourceEr{object: pod}
}

func(er *podResourceEr)PodResourceEr(...interface{})(*corev1.Pod,error){
	switch er.object.(type) {
	case *corev1.Pod:
		cm := er.object.(*corev1.Pod)
		return cm.DeepCopy(), nil
	case *crdv1alpha1.ZookeeperCluster:
		zk := er.object.(*crdv1alpha1.ZookeeperCluster)
		return newPodForZookeeper(zk, er.id), nil
	}
	return nil, fmt.Errorf("unknow type %#v", er.object)
}



func newPodForZookeeper(cluster *crdv1alpha1.ZookeeperCluster, id int)*corev1.Pod{
	podName := getSinglePodNameForZookeeperCluster(cluster,id)
	var volumeHostPathPolicy corev1.HostPathType= corev1.HostPathDirectoryOrCreate
	_ = volumeHostPathPolicy
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForZookeeperCluster(cluster)},
			Name: podName,
			Namespace: cluster.GetNamespace(),
			Annotations: getAnnotationsForPod(id),
			Labels: getLabelForZookeeperCluster(cluster),
		},
		Spec: corev1.PodSpec{
			Hostname: getPodHostNameForZookeeperCluster(cluster,id),
			Subdomain: getZookeeperClusterCommunicateServiceName(cluster),
			Containers: []corev1.Container{
				{
					Name: podName,
					Env: []corev1.EnvVar{
						{Name: "MYID", Value: strconv.Itoa(id)},
					},
					Ports: []corev1.ContainerPort{
						{Name: "client-port", ContainerPort: cluster.Spec.ClientPort},
						{Name: "leader-election", ContainerPort: cluster.Spec.SyncPort},
						{Name: "peer-port", ContainerPort: cluster.Spec.ServerPort},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name: getConfigMapNameForZookeeper(cluster),
							MountPath: "/root/zookeeper-3.4.10/conf/zoo.cfg",
							SubPath: "zoo.cfg",
						},
						{
							Name: fmt.Sprintf("%s-data", podName),
							MountPath: "/data/zookeeper/data/",
						},
						{
							Name: fmt.Sprintf("%s-log", podName),
							MountPath: "/data/zookeeper/log/",
						},
					},
					ImagePullPolicy: corev1.PullPolicy(cluster.Spec.ImagePullPolicy),
					Image: cluster.Spec.Image,
					Command: cluster.Spec.Command,
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: getConfigMapNameForZookeeper(cluster),
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{Name: getConfigMapNameForZookeeper(cluster)},
						},
					},
				},
				{
					Name: getMountPathForData(podName),
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
				{
					Name: getMountPathForLog(podName),
					VolumeSource:corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{
						},
					},
				},
			},
		},
	}
	return pod
}


func getMountPathForData(podName string)string{
	return fmt.Sprintf("%s-data", podName)
}

func getMountPathForLog(podName string)string{
	return fmt.Sprintf("%s-log", podName)
}
