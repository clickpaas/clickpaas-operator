package manager

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1lister "k8s.io/client-go/listers/core/v1"
	"l0calh0st.cn/clickpaas-operator/pkg/operator"
)

type configManager struct {
	operator.Manager
	configMapLister corev1lister.ConfigMapLister
}

func NewConfigManager(kubeClient kubernetes.Interface, configMapLister corev1lister.ConfigMapLister)*configManager{
	return &configManager{
		Manager:         operator.Manager{KubeClient: kubeClient},
		configMapLister: configMapLister,
	}
}

func (c *configManager) Create(obj interface{}, f func(interface{})(*corev1.ConfigMap,error)) (*corev1.ConfigMap, error) {
	cm,err := f(obj)
	if err != nil{
		return nil, err
	}
	return c.KubeClient.CoreV1().ConfigMaps(cm.GetNamespace()).Create(context.TODO(), cm, metav1.CreateOptions{})
}

func (c *configManager) Update(obj interface{}, f func(interface{})(*corev1.ConfigMap,error)) (*corev1.ConfigMap, error) {
	cm,err := f(obj)
	if err != nil{
		return nil, err
	}
	return c.KubeClient.CoreV1().ConfigMaps(cm.GetNamespace()).Update(context.TODO(), cm, metav1.UpdateOptions{})
}

func (c *configManager) Delete(obj interface{}, f func(interface{})(*corev1.ConfigMap,error)) error {
	cm,err := f(obj)
	if err != nil{
		return err
	}
	return c.KubeClient.CoreV1().ConfigMaps(cm.GetNamespace()).Delete(context.TODO(), cm.GetName(), metav1.DeleteOptions{})
}

func (c *configManager) Get(obj interface{}, f func(interface{})(*corev1.ConfigMap,error)) (*corev1.ConfigMap, error) {
	cm,err := f(obj)
	if err != nil{
		return nil, err
	}
	return c.configMapLister.ConfigMaps(cm.GetNamespace()).Get(cm.GetName())
}


