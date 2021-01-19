package zookeeper

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	kubeutil "l0calh0st.cn/clickpaas-operator/pkg/operator/util/kube"
)

func doOnceBootStrap(pod *corev1.Pod,kubeClient kubernetes.Interface, restConfig *rest.Config)error{
	rmGcache := []string{"/root/zookeeper-3.4.10/bin/zkCli.sh ", "rmr", "/gcache"}
	rmId := []string{"/root/zookeeper-3.4.10/bin/zkCli.sh", "rmr", "/idgenerator"}
	importId := []string{"/root/zookeeper-3.4.10/bin/zkCli.sh" ,"<" ,"/data/lib/zk_id_commands.sh"}
	importGcache := []string{"/root/zookeeper-3.4.10/bin/zkCli.sh", "<", "/data/clickconfig/zk_gcache_commands.sh"}
	_,stder,err  := kubeutil.PodExecuteWithOptions(kubeClient, restConfig, kubeutil.ExecuteOptions{
		Command:            rmGcache,
		Namespace:          pod.GetNamespace(),
		PodName:           	pod.GetName(),

	})
	if err != nil{
		return fmt.Errorf("exec rm gacahe failed %s, %s", stder, err)
	}
	_,stder,err = kubeutil.PodExecuteWithOptions(kubeClient, restConfig, kubeutil.ExecuteOptions{
		Command:            rmId,
		Namespace:          pod.GetNamespace(),
		PodName:            pod.GetName(),
	})
	if err != nil{
		return fmt.Errorf("execte rm idgenerate faile %s  %s", stder, err)
	}
	_,stder, err = kubeutil.PodExecuteWithOptions(kubeClient, restConfig, kubeutil.ExecuteOptions{
		Command:            importId,
		Namespace:          pod.GetNamespace(),
		PodName:            pod.GetName(),
	})
	if err != nil{
		return fmt.Errorf("execute import idgenerae failed, %s %s", stder, err)
	}
	_,stder,err = kubeutil.PodExecuteWithOptions(kubeClient, restConfig, kubeutil.ExecuteOptions{
		Command:            importGcache,
		Namespace:          pod.GetNamespace(),
		PodName:            pod.GetName(),
	})
	if err != nil{
		return fmt.Errorf("execute import gcache failed %s %s", stder, err)
	}
	return nil

}