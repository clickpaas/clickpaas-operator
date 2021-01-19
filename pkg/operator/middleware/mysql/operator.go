package mysql

import (
	"fmt"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	appv1lister "k8s.io/client-go/listers/apps/v1"
	corev1lister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	crdlister "l0calh0st.cn/clickpaas-operator/pkg/client/listers/middleware/v1alpha1"
	"l0calh0st.cn/clickpaas-operator/pkg/operator"
	"l0calh0st.cn/clickpaas-operator/pkg/operator/manager"
	"time"
	kubeutil "l0calh0st.cn/clickpaas-operator/pkg/operator/util/kube"
)

type mysqlOperator struct {
	kubeClient kubernetes.Interface

	mysqlClusterLister crdlister.MysqlClusterLister

	statefulSetManager operator.StatefulSetManager
	serviceManager operator.ServiceManager
	podManager operator.PodManager

}

func NewMysqlClusterOperator(
	kubeClient kubernetes.Interface,
	mysqlClusterLister crdlister.MysqlClusterLister,
	statefulSetLister appv1lister.StatefulSetLister,
	serviceLister corev1lister.ServiceLister,
	podLister corev1lister.PodLister,
	)*mysqlOperator {
	op := &mysqlOperator{
		kubeClient:         kubeClient,
		mysqlClusterLister: mysqlClusterLister,
	}
	op.statefulSetManager = manager.NewStatefulSetManager(kubeClient, statefulSetLister)
	op.serviceManager = manager.NewServiceManager(kubeClient, serviceLister)
	op.podManager = manager.NewPodManager(kubeClient, podLister)
	return op
}


func(o *mysqlOperator)Reconcile(key string)error{
	namespace,name,err := cache.SplitMetaNamespaceKey(key)
	if err != nil{
		runtime.HandleError(fmt.Errorf("invalid resource key '%s'", key))
		return nil
	}
	mc,err := o.mysqlClusterLister.MysqlClusters(namespace).Get(name)
	if err != nil{
		if k8serr.IsNotFound(err){
			runtime.HandleError(fmt.Errorf("mysql '%s:%s' is not longer existed", namespace, name))
			return nil
		}
		return fmt.Errorf("list mysql failed '%s:%s':%s", namespace, name, err)
	}
	mysqlCopy := mc.DeepCopy()
	// check statefulSet is exists, if not exited ,then create one
	mysqlSs,err := o.statefulSetManager.Get(&statefulSetEr{mysqlCopy})
	if err != nil{
		if k8serr.IsNotFound(err){
			mysqlSs,err = o.statefulSetManager.Create(&statefulSetEr{mysqlCopy})
			if err != nil{
				return err
			}
		} else {
			return err
		}
	}
	if err := operator.WaitForStatefulSetPodsReady(mysqlSs, 10 * time.Second); err!= nil{
		return err
	}
	// check service, if not existed, then create new one
	mysqlSvc,err := o.serviceManager.Get(&serviceResourceEr{mysqlCopy})
	if err != nil{
		if k8serr.IsNotFound(err){
			mysqlSvc,err = o.serviceManager.Create(&serviceResourceEr{mysqlCopy})
			if err != nil{
				return err
			}
		} else {
			return err
		}
	}
	_ = mysqlSvc
	// bootstrap all data
	allPods,err := o.podManager.List(getLabelForMysqlCluster(mysqlCopy))
	if err != nil{
		return fmt.Errorf("list mysql pod failed %s:%s  %s", mysqlCopy.GetName(), mysqlCopy.GetNamespace(), err.Error())
	}
	if err := kubeutil.WaitForPodsReady(allPods, 10 * time.Second); err != nil{
		return fmt.Errorf("wait all pod ready timeout, %s", err.Error())
	}
	if len(allPods) != len(kubeutil.FilteredActivePods(allPods)) {
		return fmt.Errorf("double check all pods actived failed, may existed some not actived pods")
	}




	return nil
}

func(o *mysqlOperator)Healthy()error{
	return nil
}