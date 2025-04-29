package bootstrap

import (
	"flag"

	conf "github.com/fzf-labs/kratos-contrib/api/conf/v1"
	v1 "github.com/fzf-labs/kratos-contrib/api/conf/v1"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

type Flags struct {
	// conf is the config flag.
	conf string
}

func NewFlags() *Flags {
	f := new(Flags)
	flag.StringVar(&f.conf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.Parse()
	return f
}

// LoadConfig 加载配置
func LoadConfig(flagconf string) *v1.Bootstrap {
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	return &bc
}
