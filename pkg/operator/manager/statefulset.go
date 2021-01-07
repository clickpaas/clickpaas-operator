package manager

import (
	"context"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appv1lister "k8s.io/client-go/listers/apps/v1"
	"l0calh0st.cn/clickpaas-operator/pkg/operator"
)

type statefulSetManager struct {
	operator.Manager
	statefulSetLister  appv1lister.StatefulSetLister
}

func NewStatefulSetManager(kubeClient kubernetes.Interface, statefulSetLister appv1lister.StatefulSetLister)*statefulSetManager{
	return &statefulSetManager{
		Manager:           operator.Manager{KubeClient: kubeClient},
		statefulSetLister: statefulSetLister,
	}
}

func (s *statefulSetManager) Create(object operator.StatefulSetResourceEr) (*appv1.StatefulSet, error) {
	ss,err := object.StatefulSetResourceEr()
	if err != nil{
		return nil, err
	}
	return s.KubeClient.AppsV1().StatefulSets(ss.GetNamespace()).Create(context.TODO(), ss, metav1.CreateOptions{})
}

func (s *statefulSetManager) Update(object operator.StatefulSetResourceEr) (*appv1.StatefulSet, error) {
	ss,err := object.StatefulSetResourceEr()
	if err != nil{
		return nil, err
	}
	return s.KubeClient.AppsV1().StatefulSets(ss.GetNamespace()).Update(context.TODO(), ss, metav1.UpdateOptions{})
}

func (s *statefulSetManager) Delete(object operator.StatefulSetResourceEr) error {
	ss,err := object.StatefulSetResourceEr()
	if err != nil{
		return err
	}
	return s.KubeClient.AppsV1().StatefulSets(ss.GetNamespace()).Delete(context.TODO(), ss.GetName(), metav1.DeleteOptions{})
}

func (s *statefulSetManager) Get(object operator.StatefulSetResourceEr) (*appv1.StatefulSet, error) {
	ss,err := object.StatefulSetResourceEr()
	if err != nil{
		return nil, err
	}
	return s.KubeClient.AppsV1().StatefulSets(ss.GetNamespace()).Get(context.TODO(), ss.GetName(), metav1.GetOptions{})
}


