package manager

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	corev1lister "k8s.io/client-go/listers/core/v1"
	"l0calh0st.cn/clickpaas-operator/pkg/operator"
)

type podManager struct {
	kubeClient kubernetes.Interface
	podLister corev1lister.PodLister
}

func (p *podManager) Create(er operator.PodResourceEr) (*corev1.Pod, error) {
	pod,err := er.PodResourceEr()
	if err != nil{
		return nil, err
	}
	return p.kubeClient.CoreV1().Pods(pod.GetNamespace()).Create(context.TODO(), pod, metav1.CreateOptions{})
}

func (p *podManager) Get(er operator.PodResourceEr) (*corev1.Pod, error) {
	pod,err := er.PodResourceEr()
	if err != nil{
		return nil, err
	}
	return p.podLister.Pods(pod.GetNamespace()).Get(pod.GetName())
}

func (p *podManager) Update(er operator.PodResourceEr) (*corev1.Pod, error) {
	pod,err := er.PodResourceEr()
	if err != nil{
		return nil, err
	}
	return p.kubeClient.CoreV1().Pods(pod.GetNamespace()).Update(context.TODO(), pod, metav1.UpdateOptions{})
}

func (p *podManager) Delete(er operator.PodResourceEr) error {
	panic("implement me")
}


func(p *podManager)List(labelsSet labels.Set)([]*corev1.Pod, error){
	return p.podLister.List(labels.SelectorFromSet(labelsSet))
}

func NewPodManager(kubeClient kubernetes.Interface, lister corev1lister.PodLister)*podManager{
	return &podManager{
		kubeClient: kubeClient,
		podLister:  lister,
	}
}
