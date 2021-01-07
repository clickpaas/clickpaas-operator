package mongo

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
)

type mongoOperator struct {
	kubeClient kubernetes.Interface
	mongoLister crdlister.MongoClusterLister

	serviceManager operator.ServiceManager
	statefulManager operator.StatefulSetManager

}


func NewMongoOperator(kubeClient kubernetes.Interface, mongoLister crdlister.MongoClusterLister,
	serviceLister corev1lister.ServiceLister, statefulSetLister appv1lister.StatefulSetLister)operator.IOperator{
	return &mongoOperator{
		kubeClient:      kubeClient,
		mongoLister:     mongoLister,
		serviceManager:  manager.NewServiceManager(kubeClient, serviceLister),
		statefulManager: manager.NewStatefulSetManager(kubeClient, statefulSetLister),
	}
}


func (m *mongoOperator) Reconcile(key string) error {
	namespace,name,err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key %v", key))
		return nil
	}
	mongo,err := m.mongoLister.MongoClusters(namespace).Get(name)
	if err != nil{
		if k8serr.IsNotFound(err){
			runtime.HandleError(fmt.Errorf("resource is not exist '%v:%v'", namespace, name))
			return nil
		} else {
			return err
		}
	}
	// check statefulset is existed
	ss,err := m.statefulManager.Get(&statefulSetResourceEr{mongo})
	if err != nil{
		if k8serr.IsNotFound(err){
			ss,err = m.statefulManager.Create(&statefulSetResourceEr{mongo})
			if err != nil{
				return err
			}
		} else {
			return err
		}
	}
	if *ss.Spec.Replicas != mongo.Spec.Replicas{
		ss.Spec.Replicas = &mongo.Spec.Replicas
		if ss,err = m.statefulManager.Update(&statefulSetResourceEr{ss});err != nil{
			return err
		}
	}
	// check service is existed
	svc,err := m.serviceManager.Get(&serviceResourceEr{mongo})
	if err != nil{
		if k8serr.IsNotFound(err){
			svc,err = m.serviceManager.Create(&serviceResourceEr{mongo})
			if err != nil{
				return err
			}
		}else {
			return err
		}
	}
	_ = svc
	return nil
}

func (m *mongoOperator) Healthy() error {
	return nil
}


