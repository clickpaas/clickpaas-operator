package manager

import (
	"context"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appv1lister "k8s.io/client-go/listers/apps/v1"
	"l0calh0st.cn/clickpaas-operator/pkg/operator"
)

type deploymentManager struct {
	operator.Manager
	deploymentLister appv1lister.DeploymentLister
}

func NewDeploymentManager(kubeClient kubernetes.Interface, deploymentLister appv1lister.DeploymentLister)*deploymentManager{
	return &deploymentManager{
		Manager:          operator.Manager{KubeClient: kubeClient},
		deploymentLister: deploymentLister,
	}
}

func (d *deploymentManager) Create(object operator.DeploymentResourceEr) (*appv1.Deployment, error) {
	deployment,err := object.DeploymentResourceEr()
	if err != nil{
		return nil,err
	}
	return d.KubeClient.AppsV1().Deployments(deployment.GetNamespace()).Create(context.TODO(),deployment, metav1.CreateOptions{})
}

func (d *deploymentManager) Update(object operator.DeploymentResourceEr) (*appv1.Deployment, error) {
	deployment,err := object.DeploymentResourceEr()
	if err != nil{
		return nil,err
	}
	return d.KubeClient.AppsV1().Deployments(deployment.GetNamespace()).Update(context.TODO(),deployment, metav1.UpdateOptions{})
}

func (d *deploymentManager) Delete(object operator.DeploymentResourceEr) error {
	deployment,err := object.DeploymentResourceEr()
	if err != nil{
		return err
	}
	return d.KubeClient.AppsV1().Deployments(deployment.GetNamespace()).Delete(context.TODO(),deployment.GetName(), metav1.DeleteOptions{})
}

func (d *deploymentManager) Get(object operator.DeploymentResourceEr) (*appv1.Deployment, error) {
	deployment,err := object.DeploymentResourceEr()
	if err != nil{
		return nil,err
	}
	return d.KubeClient.AppsV1().Deployments(deployment.GetNamespace()).Get(context.TODO(),deployment.GetName(), metav1.GetOptions{})
}



