package diamond

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

func (d *deploymentResourceEr)DeploymentResourceEr(...interface{})(*appv1.Deployment,error){
	switch d.object.(type) {
	case *appv1.Deployment:
		svc := d.object.(*appv1.Deployment)
		return svc.DeepCopy(), nil
	case *crdv1alpha1.Diamond:
		mongo := d.object.(*crdv1alpha1.Diamond)
		return newDeploymentForDiamond(mongo), nil
	}
	return nil, fmt.Errorf("unexcept type %#v", d.object)
}


func newDeploymentForDiamond(diamond *crdv1alpha1.Diamond)*appv1.Deployment{
	dp := &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForDiamond(diamond)},
			Name:            getDeploymentNameForDiamond(diamond),
			Namespace:       diamond.GetNamespace(),
		},
		Spec: appv1.DeploymentSpec{
			Replicas: &diamond.Spec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: getLabelForDiamond(diamond)},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: getLabelForDiamond(diamond)},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            getDeploymentNameForDiamond(diamond),
							Image:           diamond.Spec.Image,
							ImagePullPolicy: corev1.PullPolicy(diamond.Spec.ImagePullPolicy),
							Env: []corev1.EnvVar{
								{Name: "MYSQL_HOST", Value: diamond.Spec.Db.Host},
								{Name: "DB_USER", Value: diamond.Spec.Db.User},
								{Name: "DB_PASSWORD", Value: diamond.Spec.Db.Password},
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