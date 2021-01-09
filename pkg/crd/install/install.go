package install

import (
	apiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"l0calh0st.cn/clickpaas-operator/pkg/crd/middleware"
)

func MayAutoInstallCRDs(apiClient apiextensions.Interface)(err error){
	if err = middleware.CreateMysqlClusterCRD(apiClient); err != nil{
		return
	}
	if err = middleware.CreateDiamondCRD(apiClient);err != nil{
		return
	}
	if err = middleware.CreateMongoCRD(apiClient);err != nil{
		return
	}
	if err = middleware.CreateRedisIdGeneratorCRD(apiClient);err != nil{
		return
	}
	if err = middleware.CreateRocketmqCRD(apiClient); err != nil{
		return
	}
	if err = middleware.CreateRedisGCacheCRD(apiClient);err != nil{
		return
	}
	if err = middleware.CreateLtsJobTrackerCRD(apiClient);err != nil{
		return
	}
	if err = middleware.CreateZookeeperClusterCRD(apiClient);err != nil{
		return
	}
	return
}
