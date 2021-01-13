package rediscluster

import (
	"context"
	"fmt"
	goredis "github.com/go-redis/redis/v8"
	"strconv"
	"strings"
	"time"
)

type redisAdmin struct {
	host string
	port int
	auth string

	redisClient *goredis.Client
}

func NewRedisAdmin(host, auth string, port int)*redisAdmin{
	return &redisAdmin{
		host: host,
		port: port,
		auth: auth,
	}
}


func(adm *redisAdmin)Connect()error{
	adm.redisClient = goredis.NewClient(&goredis.Options{
		Addr: adm.host + ":" + strconv.Itoa(adm.port),
		Password: adm.auth,
	})
	ctx,cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	if err := adm.redisClient.Ping(ctx).Err(); err != nil{
		return fmt.Errorf("try connect redis %s:%d failed, %s", adm.host, adm.port, err.Error())
	}
	return nil
}


func (adm *redisAdmin)NodeId()(string,error){
	return adm.nodeId()
}


func(adm *redisAdmin)AddSlots(min, max int)error{
	if err := adm.redisClient.ClusterAddSlotsRange(context.TODO(), min, max).Err(); err != nil{
		return fmt.Errorf("add slot from %d to %d failed, %s", min, max, err.Error())
	}
	return nil
}

func(adm *redisAdmin)DisConnect()error{
	return adm.redisClient.Close()
}

func(adm *redisAdmin)nodeId()(string,error){

	nid := adm.redisClient.ClusterNodes(context.TODO())
	if nid.Err() != nil{
		return "", fmt.Errorf("get node id failed %s", nid.Err().Error())
	}
	res,err := nid.Result()
	if err != nil{
		return "", fmt.Errorf("get node id failed %s",err.Error())
	}
	nodeIdSplit := strings.Split(res, "myself")
	if len(nodeIdSplit) != 2 {
		return "", fmt.Errorf("get node id from %s faild, split failed", res)
	}
	return nodeIdSplit[0], nil
}