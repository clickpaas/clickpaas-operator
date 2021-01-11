package zookeeper

import (
	"fmt"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

const (
	ZookeeperEventReasonOnAdded = "ZookeeperOnAdded"
	ZookeeperEventReasonOnDelete = "ZookeeperOnDelete"
	ZookeeperEventReasonOnUpdate = "ZookeeperOnUpdate"

	ZookeeperEventUpdateForPodUpdate = "zkPodUpdate"
	ZookeeperEventUpdateForPodDelete = "zkPodDelete"
	ZookeeperEventUpdateForPodAdd = "zkPodAdd"
)

func eventMessage(Zookeeper *crdv1alpha1.ZookeeperCluster, eventType string)string{
	return fmt.Sprintf("Zookeeper '%s:%s' %s", Zookeeper.GetName(), Zookeeper.GetNamespace(), eventType)
}


