package diamond

import (
	"context"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appv1lister "k8s.io/client-go/listers/apps/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

type deploymentManager struct {
	kubeClient kubernetes.Interface
	deploymentLister appv1lister.DeploymentLister
}


func NewDeploymentManager(kubeClient kubernetes.Interface, deploymentLister appv1lister.DeploymentLister)*deploymentManager{
	return &deploymentManager{
		kubeClient:       kubeClient,
		deploymentLister: deploymentLister,
	}
}

func(m *deploymentManager)Create(diamond *crdv1alpha1.Diamond)(*appv1.Deployment,error){
	return m.kubeClient.AppsV1().Deployments(diamond.GetNamespace()).Create(context.TODO(), newDeploymentForDiamond(diamond), metav1.CreateOptions{})
}

func(m *deploymentManager)Update(deploy *appv1.Deployment)(*appv1.Deployment, error){
	return m.kubeClient.AppsV1().Deployments(deploy.GetNamespace()).Update(context.TODO(), deploy, metav1.UpdateOptions{})
}

func(m *deploymentManager)Get(diamond *crdv1alpha1.Diamond)(*appv1.Deployment,error){
	return m.deploymentLister.Deployments(diamond.GetNamespace()).Get(getDeploymentNameForDiamond(diamond))
}

func(m *deploymentManager)Delete(diamond *crdv1alpha1.Diamond)error{
	return m.kubeClient.AppsV1().Deployments(diamond.GetNamespace()).Delete(context.TODO(), getDeploymentNameForDiamond(diamond), metav1.DeleteOptions{})
}


func newDeploymentForDiamond(diamond *crdv1alpha1.Diamond)*appv1.Deployment{
	dp := &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForDiamond(diamond)},
			Name: getDeploymentNameForDiamond(diamond),
			Namespace: diamond.GetNamespace(),
		},
		Spec: appv1.DeploymentSpec{
			Replicas: &diamond.Spec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: getLabelForDiamond(diamond)},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: getLabelForDiamond(diamond)},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: getDeploymentNameForDiamond(diamond),
							Image: diamond.Spec.Image,
							ImagePullPolicy: corev1.PullPolicy(diamond.Spec.ImagePullPolicy),
							Env: []corev1.EnvVar{
								{Name: "MYSQL_HOST", Value: diamond.Spec.Config.Host},
								{Name: "DB_USER", Value: diamond.Spec.Config.User},
								{Name: "DB_PASSWORD", Value: diamond.Spec.Config.Password},
							},
							Ports: []corev1.ContainerPort{{Name: "diamond-port", ContainerPort: diamond.Spec.Port}},
						},
					},
				},
			},
		},
	}
	return dp
}