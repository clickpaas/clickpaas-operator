package rocketmq

import (
	"fmt"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

const (
	RocketmqEventReasonOnAdded = "RocketmqOnAdded"
	RocketmqEventReasonOnDelete = "RocketmqOnDelete"
	RocketmqEventReasonOnUpdate = "RocketmqOnUpdate"
)

func eventMessage(rocketmq *crdv1alpha1.Rocketmq, eventType string)string{
	return fmt.Sprintf("Rocketmq '%s:%s' %s", rocketmq.GetName(), rocketmq.GetNamespace(), eventType)
}

