package bootstrap

import (
	"flag"
	"fmt"
	"os"

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
			file.NewSource(fmt.Sprintf("%s/config.%s.yaml", flagconf, os.Getenv("GO_ENV"))),
		),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}
	var bc v1.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	return &bc
}
