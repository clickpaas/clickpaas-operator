package rocketmq

import (
	"fmt"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

type deploymentResourceEr struct {
	object interface{}
}

func(er *deploymentResourceEr)DeploymentResourceEr(...interface{})(*appv1.Deployment, error){
	switch er.object.(type) {
	case *appv1.Deployment:
		svc := er.object.(*appv1.Deployment)
		return svc.DeepCopy(), nil
	case *crdv1alpha1.Rocketmq:
		rocketmq := er.object.(*crdv1alpha1.Rocketmq)
		return newDeploymentForRocketmq(rocketmq), nil
	}
	return nil, fmt.Errorf("unexcept type %#v", er.object)
}

func newDeploymentForRocketmq(rocketmq *crdv1alpha1.Rocketmq)*appv1.Deployment{
	dp := &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForRocketmqCluster(rocketmq)},
			Name: getDeploymentNameForRocketmq(rocketmq),
			Namespace: rocketmq.GetNamespace(),
		},
		Spec:       appv1.DeploymentSpec{
			Replicas: &rocketmq.Spec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: getLabelForRocketmqNameServer(rocketmq)},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: getLabelForRocketmqNameServer(rocketmq)},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: getDeploymentNameForRocketmq(rocketmq),
							ImagePullPolicy: corev1.PullPolicy(rocketmq.Spec.ImagePullPolicy),
							Image: rocketmq.Spec.Image,
							Ports: []corev1.ContainerPort{
								{Name: "ns-port", ContainerPort: rocketmq.Spec.NameServerPort},
							},
						},
					},
				},
			},
		},
	}
	return dp
}