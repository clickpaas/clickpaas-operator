package idgenerator

import (
	"fmt"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	appv1lister "k8s.io/client-go/listers/apps/v1"
	corev1lister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	crdlister "l0calh0st.cn/clickpaas-operator/pkg/client/listers/middleware/v1alpha1"
	"l0calh0st.cn/clickpaas-operator/pkg/operator"
	"l0calh0st.cn/clickpaas-operator/pkg/operator/manager"
	kubeutil "l0calh0st.cn/clickpaas-operator/pkg/operator/util/kube"
)

type redisIdGeneratorOperator struct {
	kubeClient kubernetes.Interface

	idGeneratorLister crdlister.IdGenerateLister

	serviceManager operator.ServiceManager
	statefulSetManager operator.StatefulSetManager
}

func NewIdGeneratorOperator(kubeClient kubernetes.Interface, idGeneratorLister crdlister.IdGenerateLister,
	serviceLister corev1lister.ServiceLister, statefulSetLister appv1lister.StatefulSetLister)*redisIdGeneratorOperator{
	redisOperator := &redisIdGeneratorOperator{
		kubeClient:         kubeClient,
		idGeneratorLister:  idGeneratorLister,
	}
	redisOperator.statefulSetManager = manager.NewStatefulSetManager(kubeClient, statefulSetLister)
	redisOperator.serviceManager = manager.NewServiceManager(kubeClient, serviceLister)


	return redisOperator
}

func(op *redisIdGeneratorOperator)Reconcile(key string)error{
	namespace,name,err := cache.SplitMetaNamespaceKey(key)
	if err != nil{
		runtime.HandleError(fmt.Errorf("unexpect key %#v  %v", key, err))
	}
	idGenerator,err := op.idGeneratorLister.IdGenerates(namespace).Get(name)
	if err != nil{
		if k8serr.IsNotFound(err){
			runtime.HandleError(fmt.Errorf("idgenerator is not exists  '%v:%v'", name, namespace))
			return nil
		}else {
			return err
		}
	}

	allWorkerNode,err := kubeutil.GetAllWorkNode(op.kubeClient)
	if err != nil || len(allWorkerNode) <= 0{
		return fmt.Errorf("list all worknode failed: %s",err)
	}
	randomNode := allWorkerNode[rand.Intn(len(allWorkerNode))]

	// check statefulSet
	ss,err := op.statefulSetManager.Get(&statefulSetResourceEr{idGenerator, randomNode.GetName(), nil})
	if err != nil{
		if k8serr.IsNotFound(err){
			ss,err = op.statefulSetManager.Create(&statefulSetResourceEr{idGenerator, randomNode.GetName(), nil})
			if err != nil{
				return fmt.Errorf("create statefulset failed %s",  err)
			}
		}else {
			return fmt.Errorf("get statefulset failed %s", err)
		}
	}
	_ = ss
	svc,err := op.serviceManager.Get(&serviceResourceEr{idGenerator})
	if err != nil{
		if k8serr.IsNotFound(err){
			svc,err = op.serviceManager.Create(&serviceResourceEr{idGenerator})
			if err != nil{
				return fmt.Errorf("create service failed, %s", err)
			}
		}else {
			return fmt.Errorf("get service failed %s", err)
		}
	}
	_ = svc
	return nil
}
func(op *redisIdGeneratorOperator)Healthy()error{
	return nil
}

