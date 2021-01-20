package zookeeper

import (
	"fmt"
	"path"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	kubeutil "l0calh0st.cn/clickpaas-operator/pkg/operator/util/kube"
)

func doOnceBootStrap(pod *corev1.Pod, kubeClient kubernetes.Interface, restConfig *rest.Config, cluster *crdv1alpha1.ZookeeperCluster) error {
	zkCliBin := path.Join(cluster.Spec.ZkHome, "bin/zkCli.sh")
	rmGcache := []string{"sh", "-c", zkCliBin, "rmr", "/gcache"}
	rmID := []string{"sh", "-c", zkCliBin, "rmr", "/idgenerator"}
	importID := []string{"sh", "-c", zkCliBin, "<", "/tmp/lib/zk_id_commands.txt"}
	importGcache := []string{"sh", "-c", zkCliBin, "<", "/tmp/lib/zk_gcache_commands.txt"}
	_, stder, err := kubeutil.PodExecuteWithOptions(kubeClient, restConfig, kubeutil.ExecuteOptions{
		Command:       rmGcache,
		Namespace:     pod.GetNamespace(),
		PodName:       pod.GetName(),
		CaptureStdout: true,
		CaptureStderr: true,
	})
	if err != nil {
		return fmt.Errorf("exec rm gacahe failed %s, %s", stder, err)
	}
	_, stder, err = kubeutil.PodExecuteWithOptions(kubeClient, restConfig, kubeutil.ExecuteOptions{
		Command:       rmID,
		Namespace:     pod.GetNamespace(),
		PodName:       pod.GetName(),
		CaptureStdout: true,
		CaptureStderr: true,
	})
	if err != nil {
		return fmt.Errorf("execte rm idgenerate faile %s  %s", stder, err)
	}
	_, stder, err = kubeutil.PodExecuteWithOptions(kubeClient, restConfig, kubeutil.ExecuteOptions{
		Command:       importID,
		Namespace:     pod.GetNamespace(),
		PodName:       pod.GetName(),
		CaptureStdout: true,
		CaptureStderr: true,
	})
	if err != nil {
		return fmt.Errorf("execute import idgenerae failed, %s %s", stder, err)
	}
	_, stder, err = kubeutil.PodExecuteWithOptions(kubeClient, restConfig, kubeutil.ExecuteOptions{
		Command:       importGcache,
		Namespace:     pod.GetNamespace(),
		PodName:       pod.GetName(),
		CaptureStdout: true,
		CaptureStderr: true,
	})
	if err != nil {
		return fmt.Errorf("execute import gcache failed %s %s", stder, err)
	}
	return nil

}
