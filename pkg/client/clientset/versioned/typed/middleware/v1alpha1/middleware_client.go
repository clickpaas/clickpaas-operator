/*
Copyright The clickpaas-controller Authors.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	rest "k8s.io/client-go/rest"
	v1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	"l0calh0st.cn/clickpaas-operator/pkg/client/clientset/versioned/scheme"
)

type MiddlewareV1alpha1Interface interface {
	RESTClient() rest.Interface
	DiamondsGetter
	IdGeneratesGetter
	LtsJobTrackersGetter
	MongoClustersGetter
	MysqlClustersGetter
	RedisGCachesGetter
	RocketmqsGetter
	ZookeeperClustersGetter
}

// MiddlewareV1alpha1Client is used to interact with features provided by the middleware.l0calh0st.cn group.
type MiddlewareV1alpha1Client struct {
	restClient rest.Interface
}

func (c *MiddlewareV1alpha1Client) Diamonds(namespace string) DiamondInterface {
	return newDiamonds(c, namespace)
}

func (c *MiddlewareV1alpha1Client) IdGenerates(namespace string) IdGenerateInterface {
	return newIdGenerates(c, namespace)
}

func (c *MiddlewareV1alpha1Client) LtsJobTrackers(namespace string) LtsJobTrackerInterface {
	return newLtsJobTrackers(c, namespace)
}

func (c *MiddlewareV1alpha1Client) MongoClusters(namespace string) MongoClusterInterface {
	return newMongoClusters(c, namespace)
}

func (c *MiddlewareV1alpha1Client) MysqlClusters(namespace string) MysqlClusterInterface {
	return newMysqlClusters(c, namespace)
}

func (c *MiddlewareV1alpha1Client) RedisGCaches(namespace string) RedisGCacheInterface {
	return newRedisGCaches(c, namespace)
}

func (c *MiddlewareV1alpha1Client) Rocketmqs(namespace string) RocketmqInterface {
	return newRocketmqs(c, namespace)
}

func (c *MiddlewareV1alpha1Client) ZookeeperClusters(namespace string) ZookeeperClusterInterface {
	return newZookeeperClusters(c, namespace)
}

// NewForConfig creates a new MiddlewareV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*MiddlewareV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &MiddlewareV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new MiddlewareV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *MiddlewareV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new MiddlewareV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *MiddlewareV1alpha1Client {
	return &MiddlewareV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *MiddlewareV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
