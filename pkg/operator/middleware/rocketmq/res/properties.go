package res

import (
	"fmt"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

const (
	BrokerRoleSyncMaster = "SYNC_MASTER"
)

const (
	BrokerSyncModel = "ASYNC_FLUSH"
	BrokerASyncModel = "ASYNC_FLUSH"
)

const (
	BrokerSampleProperties = `brokerClusterName=%s
brokerName=%s
brokerId=0
deleteWhen=04
fileReservedTime=480
brokerRole=SYNC_MASTER
flushDiskType=ASYNC_FLUSH
storePathRootDir=/app/store-a
storePathCommitLog=/app/store-a/commitlog
messageDelayLevel=1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 40m 50m 1h 2h 6h

autoCreateTopicEnable=true
autoCreateSubscriptionGroup=true

listenPort=10910
`
)

func NewSampleBrokerProperties(rocketmq *crdv1alpha1.Rocketmq, brokerRole string)string{
	return fmt.Sprintf(BrokerSampleProperties, rocketmq.GetName(), brokerRole)
}

