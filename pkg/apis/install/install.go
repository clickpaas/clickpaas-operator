package install

import (
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/runtime"
	middlewarev1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	"l0calh0st.cn/clickpaas-operator/pkg/client/clientset/versioned/scheme"
)

func InstallGroupVersion(){
	logrus.Infof("register scheme")
	runtime.Must(middlewarev1alpha1.AddToScheme(scheme.Scheme))
}
