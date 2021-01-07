package lts

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
	case *crdv1alpha1.LtsJobTracker:
		lts := er.object.(*crdv1alpha1.LtsJobTracker)
		return newDeploymentForLtsJobTracker(lts), nil
	}
	return nil, fmt.Errorf("unexcept type %#v", er.object)
}


func newDeploymentForLtsJobTracker(lts *crdv1alpha1.LtsJobTracker)*appv1.Deployment{
	deploy := &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: getDeploymentNameForLtsJobTracker(lts),
			Namespace: lts.GetNamespace(),
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForLtsJobTracker(lts)},
		},
		Spec:       appv1.DeploymentSpec{
			Replicas: &lts.Spec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: getLabelForLtsJobTracker(lts)},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: getLabelForLtsJobTracker(lts),
				},
				Spec:       corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: getDeploymentNameForLtsJobTracker(lts),
							ImagePullPolicy: corev1.PullPolicy(lts.Spec.ImagePullPolicy),
							Image: lts.Spec.Image,
							Env: []corev1.EnvVar{
								{Name: "REGISTRY_ADDRESS", Value: lts.Spec.Config.RegistryAddress},
								{Name: "MYSQL_HOST", Value: lts.Spec.Config.Db.Host},
								{Name: "DB_USERNAME", Value: lts.Spec.Config.Db.User},
								{Name: "DB_PASSWORD", Value: lts.Spec.Config.Db.Password},
							},
						},
					},
				},
			},
		},
		Status:     appv1.DeploymentStatus{},
	}
	return deploy
}
