package bootstrap

import (
	conf "github.com/fzf-labs/kratos-contrib/api/conf/v1"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
)

// Bootstrap 应用引导启动
func Bootstrap(service *Service) (*conf.Bootstrap, log.Logger, registry.Registrar, registry.Discovery) {
	// new flags
	flags := NewFlags()
	// load configs
	cfg := LoadConfig(flags.conf)
	if cfg == nil {
		panic("load config failed")
	}
	// init logger
	ll := NewLoggerProvider(cfg.Logger, service)
	// init registrar
	reg, dis := NewRegistryAndDiscovery(cfg.Registry)
	// init tracer
	err := NewTracerProvider(cfg.Trace, service)
	if err != nil {
		panic(err)
	}
	return cfg, ll, reg, dis
}
