package lts

import (
	"fmt"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

const (
	LtsJobTrackerEventReasonOnAdded = "LtsJobTrackerOnAdded"
	LtsJobTrackerEventReasonOnDelete = "LtsJobTrackerOnDelete"
	LtsJobTrackerEventReasonOnUpdate = "LtsJobTrackerOnUpdate"
)

func eventInfoMessage(lts *crdv1alpha1.LtsJobTracker, eventType string)string{
	return fmt.Sprintf("LtsJobTracker '%s:%s' %s", lts.GetName(), lts.GetNamespace(), eventType)
}


func eventErrMessage(lts *crdv1alpha1.LtsJobTracker, eventType string, err error)string{
	return fmt.Sprintf("LtsjonTracker '%s:%s' EventType: %s  Err: %s",lts.GetName(), lts.GetNamespace(), eventType, err.Error())
}