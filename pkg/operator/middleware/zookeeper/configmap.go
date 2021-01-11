package zookeeper

import (
	"bytes"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	"text/template"
)

var (
	ZookeeperConfigTpl = `initLimit=100
syncLimit=50
dataDir=/data/zookeeper/data
dataLogDir=/data/zookeeper/log
clientPort=2181
autopurge.snapRetainCount=10
autopurge.purgeInterval=10

{{range .ServerList}}
{{. -}}
{{end}}

`
)

const (
	ZooKeeperConfigName = "zoo.cfg"
)

type configMapResourceEr struct {
	object interface{}
	f func(cluster *crdv1alpha1.ZookeeperCluster)(*corev1.ConfigMap,error)
}


func(er *configMapResourceEr)ConfigMapResourceEr(obj ...interface{})(*corev1.ConfigMap,error){
	switch er.object.(type) {
	case *corev1.ConfigMap:
		cm := er.object.(*corev1.ConfigMap)
		return cm.DeepCopy(), nil
	case *crdv1alpha1.ZookeeperCluster:
		zk := er.object.(*crdv1alpha1.ZookeeperCluster)
		return er.f(zk)
	}
	return nil, fmt.Errorf("unknow type %#v", er.object)
}

func newConfigMapForZookeeper(zookeeper *crdv1alpha1.ZookeeperCluster)(*corev1.ConfigMap,error){
	zkCfgData,err := generateZookeeperConfig(zookeeper)
	if err != nil{
		return nil, fmt.Errorf("generate zookeer config failed %s", err)
	}
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerReferenceForZookeeperCluster(zookeeper)},
			Name: getConfigMapNameForZookeeper(zookeeper),
			Namespace: zookeeper.GetNamespace(),
		},
		Data: map[string]string{ZooKeeperConfigName: zkCfgData},
		BinaryData: nil,
	}
	return cm, nil
}


func generateZookeeperConfig(cluster *crdv1alpha1.ZookeeperCluster)(string,error){
	tpl := template.New("configmap")
	t,err := tpl.Parse(ZookeeperConfigTpl)
	if err!= nil{
		return "", err
	}
	var outString bytes.Buffer
	ssName := getZookeeperClusterCommunicateServiceName(cluster)
	payload := struct {
		ServerList []string
	}{ServerList: []string{}}
	for i := 0 ; i < int(cluster.Spec.Replicas); i++{
		// podHost.podSubDomain.namespace.
		server := fmt.Sprintf("server.%d=%s.%s.%s:%d:%d", i, getPodHostNameForZookeeperCluster(cluster, i),
			ssName, cluster.GetNamespace(),
			cluster.Spec.ServerPort, cluster.Spec.SyncPort)
		payload.ServerList = append(payload.ServerList, server)
	}

	err = t.Execute(&outString, payload)
	if err != nil{
		return "", err
	}
	return outString.String(), nil
}

