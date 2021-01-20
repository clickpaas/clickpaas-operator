package kube

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetAllWorkNode(kubeClient kubernetes.Interface)([]*corev1.Node, error){
	nodes, err := kubeClient.CoreV1().Nodes().List(context.TODO(),metav1.ListOptions{})
	if err != nil{
		return nil, err
	}
	nodeList := []*corev1.Node{}
	for _, node := range nodes.Items{
		nodeList = append(nodeList, &node)
	}
	return nodeList, nil
}
