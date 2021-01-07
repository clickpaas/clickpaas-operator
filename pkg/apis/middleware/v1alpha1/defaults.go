package v1alpha1


const (
	defaultImagePullPolicy = "IfNotPresent"
)

// WithDefaultsMysqlCluster default not specified fields
func WithDefaultsMysqlCluster(cluster *MysqlCluster){
	if len(cluster.Spec.Args) == 0{
		cluster.Spec.Args = []string{"--character-set-server=utf8mb4", "--collation-server=utf8mb4_unicode_ci", "--lower_case_table_names=1"}
	}
	if cluster.Spec.Config.Password == "" {
		cluster.Spec.Config.Password = "diamond^^^"
	}
	if cluster.Spec.Config.User == ""{
		cluster.Spec.Config.User = "root"
	}
	if cluster.Spec.ImagePullPolicy == ""{
		cluster.Spec.ImagePullPolicy = defaultImagePullPolicy
	}
	if cluster.Spec.Replicas == 0{
		cluster.Spec.Replicas = 1
	}
	if cluster.Spec.Image == ""{
		cluster.Spec.Image = "registry.bizsaas.net/mysql:5.7.22"
	}
	if cluster.Spec.Port == 0{
		cluster.Spec.Port = 3306
	}
	// no
}


// WithDefaultsDiamond set default value for not special fields
func WithDefaultsDiamond(diamond *Diamond){
	if diamond.Spec.Port == 0{
		diamond.Spec.Port = 80
	}
	if diamond.Spec.Db.User == ""{
		diamond.Spec.Db.User = "root"
	}
	if diamond.Spec.Db.Password == ""{
		diamond.Spec.Db.Password = "diamond^^^"
	}
	if diamond.Spec.Replicas == 0{
		diamond.Spec.Replicas = 1
	}
	if diamond.Spec.ImagePullPolicy == ""{
		diamond.Spec.ImagePullPolicy = defaultImagePullPolicy
	}
	if diamond.Spec.Image == ""{
		diamond.Spec.Image = "registry.bizsaas.net/diamond:2.0.0-r2"
	}
}


func WithDefaultsMongoCluster(cluster *MongoCluster){
	if cluster.Spec.Port == 0 {
		cluster.Spec.Port = 27017
	}
	if cluster.Spec.Image == ""{
		cluster.Spec.Image = "registry.bizsaas.net/mongo:4.2.8"
	}
	if cluster.Spec.ImagePullPolicy == ""{
		cluster.Spec.ImagePullPolicy = defaultImagePullPolicy
	}
	if cluster.Spec.Config.User == ""{
		cluster.Spec.Config.User = "admin"
	}
	if cluster.Spec.Config.Password == ""{
		cluster.Spec.Config.Password = "admin"
	}
	if cluster.Spec.Replicas == 0{
		cluster.Spec.Replicas = 1
	}
}


func WithDefaultsRedisGCache(redisCache *RedisGCache){
	if redisCache.Spec.Port == 0{
		redisCache.Spec.Port = 8300
	}
	if redisCache.Spec.Image == ""{
		redisCache.Spec.Image = "registry.bizsaas.net/redis:cluster-singleton-3.0.7"
	}
	if redisCache.Spec.ImagePullPolicy == ""{
		redisCache.Spec.ImagePullPolicy = defaultImagePullPolicy
	}
	if redisCache.Spec.Replicas == 0{
		redisCache.Spec.Replicas = 1
	}


}

func WithDefaultsRedisIdGenerate(generate *IdGenerate){
	if generate.Spec.Port == 0{
		generate.Spec.Port = 16379
	}
	if generate.Spec.ImagePullPolicy == ""{
		generate.Spec.ImagePullPolicy = defaultImagePullPolicy
	}
	if generate.Spec.Image == ""{
		generate.Spec.Image = "registry.bizsaas.net/redis:sync-3.0.7"
	}
	if generate.Spec.Replicas == 0{
		generate.Spec.Replicas = 1
	}
}


func WithDefaultsRocketmq(rocketmq *Rocketmq){
	if rocketmq.Spec.Image == ""{
		rocketmq.Spec.Image = "registry.bizsaas.net/rocketmq:3.2.6-r1"
	}
	if rocketmq.Spec.ImagePullPolicy == ""{
		rocketmq.Spec.ImagePullPolicy = defaultImagePullPolicy
	}
	if rocketmq.Spec.Replicas == 0{
		rocketmq.Spec.Replicas = 1
	}
	if rocketmq.Spec.NameServerPort == 0{
		rocketmq.Spec.NameServerPort = 9876
	}
	if rocketmq.Spec.FastPort == 0{
		rocketmq.Spec.FastPort = 10910
	}
	if rocketmq.Spec.ListenPort == 0{
		rocketmq.Spec.ListenPort = 10911
	}
	if rocketmq.Spec.HaPort == 0{
		rocketmq.Spec.HaPort = 10912
	}
}


func WithDefaultsLtsJobTracker(lts *LtsJobTracker){
	if lts.Spec.ImagePullPolicy == ""{
		lts.Spec.ImagePullPolicy = defaultImagePullPolicy
	}
	if lts.Spec.Image == ""{
		lts.Spec.Image = "registry.bizsaas.net/arch/lts:1.6.9-r1"
	}
	if lts.Spec.HealthPort == 0{
		lts.Spec.HealthPort = 3502
	}
	if lts.Spec.Config.RegistryAddress == ""{
		lts.Spec.Config.RegistryAddress = "zookeeper://zookeeper0-0.zookeeper.default.svc.cluster.local:2181,zookeeper1-0.zookeeper.default.svc.cluster.local:2181,zookeeper2-0.zookeeper.default.svc.cluster.local:2181"
	}
	if lts.Spec.Config.Db.Host == ""{
		lts.Spec.Config.Db.Host = "diamond-mysql.default.svc.cluster.local"
	}
	if lts.Spec.Config.Db.User == ""{
		lts.Spec.Config.Db.User = "root"
	}
	if lts.Spec.Config.Db.Password == ""{
		lts.Spec.Config.Db.Password = "diamond^^^"
	}

}