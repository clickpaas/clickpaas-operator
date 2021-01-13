package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	install2 "l0calh0st.cn/clickpaas-operator/pkg/apis/install"
	crdclient "l0calh0st.cn/clickpaas-operator/pkg/client/clientset/versioned"
	crdinformer "l0calh0st.cn/clickpaas-operator/pkg/client/informers/externalversions"
	"l0calh0st.cn/clickpaas-operator/pkg/controller"
	"l0calh0st.cn/clickpaas-operator/pkg/controller/middleware/diamond"
	"l0calh0st.cn/clickpaas-operator/pkg/controller/middleware/gcache"
	"l0calh0st.cn/clickpaas-operator/pkg/controller/middleware/idgenerate"
	"l0calh0st.cn/clickpaas-operator/pkg/controller/middleware/lts"
	"l0calh0st.cn/clickpaas-operator/pkg/controller/middleware/mongo"
	"l0calh0st.cn/clickpaas-operator/pkg/controller/middleware/mysql"
	"l0calh0st.cn/clickpaas-operator/pkg/controller/middleware/rocketmq"
	"l0calh0st.cn/clickpaas-operator/pkg/controller/middleware/zookeeper"
	"l0calh0st.cn/clickpaas-operator/pkg/crd/install"

)

var (
	masterUrl = pflag.String("masterUrl", "", "address of k8s apiServer")
	kubeConfig = pflag.String("kubeConfig", "", "path of kubeConfig")
	workThreads = pflag.Int("workThreads", 1, "the number of work threads")
	resyncInterval = pflag.Int("resyncInterval", 5, "resyncInterval")
	// special options
	namespace = pflag.String("namespace", "", "the namespaced scope")
	labelSelectorFilter = pflag.String("labelSelectFilter", "", "special label filter")
)

var (
	kubeClient kubernetes.Interface
	crdClient crdclient.Interface
	extClient apiextensions.Interface
	restConfig *rest.Config
)

func init(){

}

func main(){
	pflag.Parse()

	if err := buildKubeConfig(*masterUrl, *kubeConfig); err != nil{
		logrus.Panicf("create k8s config failed, %s",err)
	}
	if err := buildKubeAndCrdClients(restConfig); err != nil{
		logrus.Panicf("build kubernetes client and crd client failed, %s", err)
	}
	ctx,cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := registerCrd();err != nil{
		logrus.Panicf("create crd failed, %s", err)
	}
	registerGVK()

	crdInformer := buildCrdInformerFactory(crdClient)
	kubeInformer := buildStandardInformerFactory(kubeClient)


	mysqlController := mysql.NewMysqlController(kubeClient, crdClient, crdInformer, kubeInformer)
	diamondController := diamond.NewDiamondController(kubeClient, crdClient, kubeInformer, crdInformer)
	mongoController := mongo.NewMongoController(kubeClient, crdClient, kubeInformer, crdInformer)
	redisGCacheController := gcache.NewRedisGCacheController(kubeClient, crdClient, kubeInformer, crdInformer)
	redisIdGenerateController := idgenerate.NewRedisIdGeneratorController(kubeClient, crdClient, kubeInformer, crdInformer)
	rocketmqController := rocketmq.NewRocketmqController(kubeClient, crdClient, kubeInformer, crdInformer)
	zookeeperController := zookeeper.NewZookeeperController(kubeClient, crdClient, crdInformer,kubeInformer)
	ltsController := lts.NewLtsJobTrackerController(kubeClient, crdClient, kubeInformer, crdInformer)

	go crdInformer.Start(ctx.Done())
	go kubeInformer.Start(ctx.Done())

	go runController(ctx, mysqlController)
	go runController(ctx,diamondController)
	go runController(ctx, mongoController)
	go runController(ctx, redisIdGenerateController)
	go runController(ctx, rocketmqController)
	go runController(ctx, zookeeperController)
	go runController(ctx, ltsController)
	go runController(ctx, redisGCacheController)


	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGABRT)

	select {
	case <-stopCh:
		os.Exit(-1)
	}
}

func runController(ctx context.Context, controller controller.IController){
	if workThreads == nil || *workThreads == 0{
		threadNumber := 1
		workThreads = &threadNumber
	}
	logrus.Infof("ready to run controller with %d threads", *workThreads)
	if err := controller.Start(ctx, *workThreads); err != nil{
		logrus.Panic("Start Controller failed, %s", err)
	}
}


func buildKubeConfig(masterUrl, kubeConfig string)(err error){
	if kubeConfig != ""{
		restConfig,err = clientcmd.BuildConfigFromFlags(masterUrl, kubeConfig)
	} else {
		restConfig,err = rest.InClusterConfig()
	}
	return
}

func buildKubeAndCrdClients(restConfig *rest.Config)(err error){
	if restConfig == nil{
		return fmt.Errorf("*rest.Config is nill")
	}
	if kubeClient,err = kubernetes.NewForConfig(restConfig);err != nil{
		return
	}
	if crdClient,err = crdclient.NewForConfig(restConfig);err != nil{
		return
	}
	if extClient,err = apiextensions.NewForConfig(restConfig); err != nil{
		return
	}
	return
}

func buildCrdInformerFactory(crdClient crdclient.Interface)crdinformer.SharedInformerFactory{
	var factoryOpts []crdinformer.SharedInformerOption
	if *namespace != corev1.NamespaceAll {
		factoryOpts = append(factoryOpts, crdinformer.WithNamespace(*namespace))
	}
	if len(*labelSelectorFilter) > 0 {
		tweakListOptionsFunc := func(options *metav1.ListOptions) {
			options.LabelSelector = *labelSelectorFilter
		}
		factoryOpts = append(factoryOpts, crdinformer.WithTweakListOptions(tweakListOptionsFunc))
	}
	return crdinformer.NewSharedInformerFactoryWithOptions(crdClient, time.Duration(*resyncInterval)*time.Second, factoryOpts...)
}

func buildStandardInformerFactory(kubeClient kubernetes.Interface)informers.SharedInformerFactory{
	var factoryOpts []informers.SharedInformerOption
	if *namespace != corev1.NamespaceAll{
		factoryOpts = append(factoryOpts, informers.WithNamespace(*namespace))
	}
	if len(*labelSelectorFilter) > 0{
		tweakListOptionsFunc := func(options *metav1.ListOptions) {
			options.LabelSelector = *labelSelectorFilter
		}
		factoryOpts = append(factoryOpts, informers.WithTweakListOptions(tweakListOptionsFunc))
	}
	return informers.NewSharedInformerFactoryWithOptions(kubeClient, time.Duration(*resyncInterval) * time.Second, factoryOpts...)
}

func registerGVK(){
	install2.InstallGroupVersion()
}

func registerCrd()(err error){
	// create mysql
	if err = install.MayAutoInstallCRDs(extClient);err != nil{
		return
	}
	return
}