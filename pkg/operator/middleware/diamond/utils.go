package diamond

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

func getDeploymentNameForDiamond(diamond *crdv1alpha1.Diamond)string{
	return fmt.Sprintf("%s", diamond.GetName())
}

func getConfigMapNameForDiamond(diamond *crdv1alpha1.Diamond)string{
	return fmt.Sprintf("%s", diamond.GetName())
}

func getServiceNameForDiamond(diamond *crdv1alpha1.Diamond)string{
	return fmt.Sprintf("%s", diamond.GetName())
}


func getLabelForDiamond(diamond *crdv1alpha1.Diamond)map[string]string{
	return map[string]string{"crdversion": crdv1alpha1.MiddlewareResourceVersion,
		"appname": diamond.GetName(),
		"kind": crdv1alpha1.DiamondKind,
	}
}

func ownerReferenceForDiamond(diamond *crdv1alpha1.Diamond)metav1.OwnerReference{
	return *metav1.NewControllerRef(diamond, crdv1alpha1.SchemeGroupVersion.WithKind(crdv1alpha1.DiamondKind))
}