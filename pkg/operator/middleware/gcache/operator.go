package gcache

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
	"l0calh0st.cn/clickpaas-operator/pkg/operator/middleware/gcache/rediscluster"
)

type redisGCacheOperator struct {
	kubeClient kubernetes.Interface

	gCacheLister crdlister.RedisGCacheLister

	statefulSetManager operator.StatefulSetManager
	serviceManager operator.ServiceManager
}



func NewRedisGCacheOperator(kubeClient kubernetes.Interface, redisGCacheLister crdlister.RedisGCacheLister, statefulSetLister appv1lister.StatefulSetLister,
	serviceLister corev1lister.ServiceLister)operator.IOperator{
	op := &redisGCacheOperator{
		kubeClient:         kubeClient,
		gCacheLister:       redisGCacheLister,
	}
	op.statefulSetManager = manager.NewStatefulSetManager(kubeClient, statefulSetLister)
	op.serviceManager = manager.NewServiceManager(kubeClient, serviceLister)

	return op
}


func (op *redisGCacheOperator) Reconcile(key string) error {
	namespace, name ,err := cache.SplitMetaNamespaceKey(key)
	if err != nil{
		runtime.HandleError(fmt.Errorf("unexcepted key %v", key))
		return nil
	}
	redisGCache,err := op.gCacheLister.RedisGCaches(namespace).Get(name)
	if err != nil{
		if k8serr.IsNotFound(err){
			runtime.HandleError(fmt.Errorf("resource redisGCache is not existed in workqueue"))
			return nil
		}else {
			return err
		}
	}
	// check statefulset
	ss ,err := op.statefulSetManager.Get(&statefulSetResourceEr{redisGCache})
	if err != nil{
		if k8serr.IsNotFound(err){
			ss,err = op.statefulSetManager.Create(&statefulSetResourceEr{redisGCache})
			if err != nil{
				return err
			}
		}else {
			return err
		}
	}
	_ = ss
	// todo check statefulset
	svc,err := op.serviceManager.Get(&serviceResourceEr{redisGCache})
	if err != nil{
		if k8serr.IsNotFound(err){
			svc,err = op.serviceManager.Create(&serviceResourceEr{redisGCache})
			if err != nil{
				return err
			}
		}else {
			return err
		}
	}
	_ = svc
	redisAdm := rediscluster.NewRedisAdmin(getServiceNameForRedisGCache(redisGCache), "", int(redisGCache.Spec.Port))
	if err := redisAdm.Connect();err != nil{
		return err
	}
	defer redisAdm.DisConnect()
	if err := redisAdm.AddSlots(0, 16383);err != nil{
		return err
	}
	return nil
}

func (op *redisGCacheOperator) Healthy() error {
	panic("implement me")
}