package rocketmq

import (
	"fmt"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

type statefulSetResourceEr struct {
	object interface{}
}


func(er *statefulSetResourceEr)StatefulSetResourceEr(...interface{})(*appv1.StatefulSet,error){
	switch er.object.(type) {
	case *appv1.StatefulSet:
		svc := er.object.(*appv1.StatefulSet)
		return svc.DeepCopy(), nil
	case *crdv1alpha1.Rocketmq:
		rocketmq := er.object.(*crdv1alpha1.Rocketmq)
		return newStatefulSetForRocketmq(rocketmq), nil
	}
	return nil, fmt.Errorf("unexcept type %#v", er.object)
}

func newStatefulSetForRocketmq(rocketmq *crdv1alpha1.Rocketmq)*appv1.StatefulSet{
	var customeCommand []string
	if len(rocketmq.Spec.Command) == 0{
		customeCommand = []string{"sh", "/app/alibaba-rocketmq-20150824/bin/mqnamesrv", "-c", "/opt/"+getBrokerPropertiesFileName(rocketmq)}
	}else {
		customeCommand = []string{}
	}
	ss := &appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForRocketmqCluster(rocketmq)},
			Name: getStatefulSetNameForRocketmq(rocketmq),
			Namespace: rocketmq.GetNamespace(),
		},
		Spec:       appv1.StatefulSetSpec{
			//Replicas: &rocketmq.Spec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: getLabelForRocketmqCluster(rocketmq)},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: getLabelForRocketmqCluster(rocketmq),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: getStatefulSetNameForRocketmq(rocketmq),
							Image: rocketmq.Spec.Image,
							ImagePullPolicy: corev1.PullPolicy(rocketmq.Spec.ImagePullPolicy),
							Ports: []corev1.ContainerPort{
								{Name: "haport", ContainerPort: rocketmq.Spec.HaPort},
								{Name: "listeen", ContainerPort: rocketmq.Spec.ListenPort},
								{Name: "fastport", ContainerPort: rocketmq.Spec.FastPort},
							},
							Env: []corev1.EnvVar{
								{Name: "JAVA_HOME", Value: "/usr/lib/jvm/java-1.8-openjdk"},
								{Name: "ROCKETMQ_HOME", Value: "/app/alibaba-rocketmq-20150824"},
								{Name: "NAMESRV_ADDR", Value: getServiceNameForRocketNameServer(rocketmq)},
							},
							Command: append(customeCommand, "-n", getServiceNameForRocketNameServer(rocketmq)),
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: getVolumeNameForBrokerProperties(rocketmq),
									MountPath: "/opt/" + getBrokerPropertiesFileName(rocketmq),
									SubPath: getBrokerPropertiesFileName(rocketmq),
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: getVolumeNameForBrokerProperties(rocketmq),
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{getConfigMapNameForBrokerProperties(rocketmq)},
								},
							},
						},
					},
				},
			},
		},
	}
	return ss
}