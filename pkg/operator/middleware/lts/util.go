package lts

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	"fmt"
)

func getDeploymentNameForLtsJobTracker(lts *crdv1alpha1.LtsJobTracker)string{
	return fmt.Sprintf("%s-lts", lts.GetName())
}

func getConfigMapNameForLtsJobTracker(lts *crdv1alpha1.LtsJobTracker)string{
	return fmt.Sprintf("%s-lts", lts.GetName())
}

func getServiceNameForLtsJobTracker(lts *crdv1alpha1.LtsJobTracker)string{
	return fmt.Sprintf("%s-lts", lts.GetName())
}


func getLabelForLtsJobTracker(lts *crdv1alpha1.LtsJobTracker)map[string]string{
	return map[string]string{"crdversion": crdv1alpha1.MiddlewareResourceVersion,
		"appname": lts.GetName(),
		"kind": crdv1alpha1.LtsJobTrackerKind,
	}
}

func ownerReferenceForLtsJobTracker(lts *crdv1alpha1.LtsJobTracker)metav1.OwnerReference{
	return *metav1.NewControllerRef(lts, crdv1alpha1.SchemeGroupVersion.WithKind(crdv1alpha1.LtsJobTrackerKind))
}
