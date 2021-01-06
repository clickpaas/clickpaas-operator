package mysql

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appv1lister "k8s.io/client-go/listers/apps/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	appv1 "k8s.io/api/apps/v1"
)

type statefulSetManager struct {
	kubeClient kubernetes.Interface
	statefulSetLister appv1lister.StatefulSetLister
}

func NewStatefulSetManager(kubeClient kubernetes.Interface, ssLister appv1lister.StatefulSetLister)*statefulSetManager{
	return &statefulSetManager{
		kubeClient:        kubeClient,
		statefulSetLister: ssLister,
	}
}

func (m *statefulSetManager)Create(cluster *crdv1alpha1.MysqlCluster)(*appv1.StatefulSet,error){
	ss := newStatefulSetForMysqlCluster(cluster)
	return m.kubeClient.AppsV1().StatefulSets(ss.GetNamespace()).Create(context.TODO(), ss, metav1.CreateOptions{})
}

func (m *statefulSetManager)Delete(cluster *crdv1alpha1.MysqlCluster)error{
	return m.kubeClient.AppsV1().StatefulSets(cluster.GetNamespace()).Delete(context.TODO(), getStatefulSetNameForMysql(cluster), metav1.DeleteOptions{})
}

func (m *statefulSetManager)Update(ss *appv1.StatefulSet)(*appv1.StatefulSet,error){
	return m.kubeClient.AppsV1().StatefulSets(ss.GetNamespace()).Update(context.TODO(), ss, metav1.UpdateOptions{})
}

func (m *statefulSetManager)Get(cluster *crdv1alpha1.MysqlCluster)(*appv1.StatefulSet, error){
	ss,err := m.statefulSetLister.StatefulSets(cluster.GetNamespace()).Get(cluster.GetName())
	return ss, err
}


func newStatefulSetForMysqlCluster(cluster *crdv1alpha1.MysqlCluster)*appv1.StatefulSet{
	ss := &appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: getStatefulSetNameForMysql(cluster),
			Namespace: cluster.GetNamespace(),
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForMysqlCluster(cluster)},
		},
		Spec: appv1.StatefulSetSpec{
			Replicas: &cluster.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: getLabelForMysqlCluster(cluster),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: getLabelForMysqlCluster(cluster)},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyAlways,
					Containers: []corev1.Container{
						{
							ImagePullPolicy: corev1.PullPolicy(cluster.Spec.ImagePullPolicy),
							Env: []corev1.EnvVar{
								{Name: "MYSQL_ROOT_PASSWORD", Value: cluster.Spec.Config.Password},
								{Name: "MYSQL_ROOT_HOST", Value: "%"},
							},
							Args: cluster.Spec.Args,
							Ports: []corev1.ContainerPort{
								{Name: "mysql-port", ContainerPort: cluster.Spec.Port},
							},
							Image: cluster.Spec.Image,
							Name: getStatefulSetNameForMysql(cluster),
						},
					},
				},
			},
		},
	}
	return ss
}