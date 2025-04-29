package bootstrap

import (
	"path/filepath"

	conf "github.com/fzf-labs/kratos-contrib/api/conf/v1"
	consulKratos "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	etcdKratos "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	k8sRegistry "github.com/go-kratos/kratos/contrib/registry/kubernetes/v2"
	nacosKratos "github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	consulClient "github.com/hashicorp/consul/api"
	nacosClients "github.com/nacos-group/nacos-sdk-go/clients"
	nacosConstant "github.com/nacos-group/nacos-sdk-go/common/constant"
	nacosVo "github.com/nacos-group/nacos-sdk-go/vo"
	etcdClient "go.etcd.io/etcd/client/v3"
	k8s "k8s.io/client-go/kubernetes"
	k8sRest "k8s.io/client-go/rest"
	k8sTools "k8s.io/client-go/tools/clientcmd"
	k8sUtil "k8s.io/client-go/util/homedir"
)

type RegistryType string

const (
	Consul     RegistryType = "consul"
	Etcd       RegistryType = "etcd"
	Nacos      RegistryType = "nacos"
	Kubernetes RegistryType = "kubernetes"
)

// NewRegistryAndDiscovery 创建一个服务发现客户端
func NewRegistryAndDiscovery(cfg *conf.Registry) (registry.Registrar, registry.Discovery) {
	if cfg == nil {
		return nil, nil
	}
	switch RegistryType(cfg.Type) {
	case Consul:
		res := NewConsulRegistry(cfg)
		return res, res
	case Etcd:
		res := NewEtcdRegistry(cfg)
		return res, res
	case Nacos:
		res := NewNacosRegistry(cfg)
		return res, res
	case Kubernetes:
		res := NewKubernetesRegistry(cfg)
		return res, res
	}
	return nil, nil
}

func NewRegistry(cfg *conf.Registry) registry.Registrar {
	if cfg == nil {
		return nil
	}
	switch RegistryType(cfg.Type) {
	case Consul:
		res := NewConsulRegistry(cfg)
		return res
	case Etcd:
		res := NewEtcdRegistry(cfg)
		return res
	case Nacos:
		res := NewNacosRegistry(cfg)
		return res
	case Kubernetes:
		res := NewKubernetesRegistry(cfg)
		return res
	}
	return nil
}

// NewConsulRegistry 创建一个注册发现客户端 - Consul
func NewConsulRegistry(c *conf.Registry) *consulKratos.Registry {
	cfg := consulClient.DefaultConfig()
	cfg.Address = c.Consul.Address
	cfg.Scheme = c.Consul.Scheme

	var cli *consulClient.Client
	var err error
	if cli, err = consulClient.NewClient(cfg); err != nil {
		log.Fatal(err)
	}

	reg := consulKratos.New(cli, consulKratos.WithHealthCheck(c.Consul.HealthCheck))

	return reg
}

// NewEtcdRegistry 创建一个注册发现客户端 - Etcd
func NewEtcdRegistry(c *conf.Registry) *etcdKratos.Registry {
	cfg := etcdClient.Config{
		Endpoints: c.Etcd.Endpoints,
	}

	var err error
	var cli *etcdClient.Client
	if cli, err = etcdClient.New(cfg); err != nil {
		log.Fatal(err)
	}

	reg := etcdKratos.New(cli)

	return reg
}

// NewNacosRegistry 创建一个注册发现客户端 - Nacos
func NewNacosRegistry(c *conf.Registry) *nacosKratos.Registry {
	srvConf := []nacosConstant.ServerConfig{
		*nacosConstant.NewServerConfig(c.Nacos.Address, c.Nacos.Port),
	}

	cliConf := nacosConstant.ClientConfig{
		NamespaceId:          c.Nacos.NamespaceId,
		TimeoutMs:            uint64(c.Nacos.Timeout.AsDuration().Milliseconds()), // http请求超时时间，单位毫秒
		BeatInterval:         c.Nacos.BeatInterval.AsDuration().Milliseconds(),    // 心跳间隔时间，单位毫秒
		UpdateThreadNum:      int(c.Nacos.UpdateThreadNum),                        // 更新服务的线程数
		LogLevel:             c.Nacos.LogLevel,
		CacheDir:             c.Nacos.CacheDir,             // 缓存目录
		LogDir:               c.Nacos.LogDir,               // 日志目录
		NotLoadCacheAtStart:  c.Nacos.NotLoadCacheAtStart,  // 在启动时不读取本地缓存数据，true--不读取，false--读取
		UpdateCacheWhenEmpty: c.Nacos.UpdateCacheWhenEmpty, // 当服务列表为空时是否更新本地缓存，true--更新,false--不更新
	}

	cli, err := nacosClients.NewNamingClient(
		nacosVo.NacosClientParam{
			ClientConfig:  &cliConf,
			ServerConfigs: srvConf,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	reg := nacosKratos.New(cli)

	return reg
}

// NewKubernetesRegistry 创建一个注册发现客户端 - Kubernetes
func NewKubernetesRegistry(_ *conf.Registry) *k8sRegistry.Registry {
	restConfig, err := k8sRest.InClusterConfig()
	if err != nil {
		home := k8sUtil.HomeDir()
		kubeConfig := filepath.Join(home, ".kube", "config")
		restConfig, err = k8sTools.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			log.Fatal(err)
			return nil
		}
	}

	clientSet, err := k8s.NewForConfig(restConfig)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	reg := k8sRegistry.NewRegistry(clientSet)

	return reg
}
