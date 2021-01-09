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
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForZookeeperCluster(cluster)},
			Name: getSinglePodNameForZookeeperCluster(cluster, id),
			Namespace: cluster.GetNamespace(),
			Annotations: getAnnotationsForPod(id),
		},
		Spec: corev1.PodSpec{
			Hostname: getPodHostNameForZookeeperCluster(cluster,id),
			Subdomain: getPodSubDomainForZookeeperCluster(cluster),
			Containers: []corev1.Container{
				{
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
					},
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
			},
		},
	}
	return pod
}