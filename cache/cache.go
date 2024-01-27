package cache

import (
	conf "fkratos-contrib/api/conf/v1"
	"fmt"

	"github.com/fzf-labs/fpkg/cache/rueidiscache"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/rueidis"
)

// NewRueidis 初始化Rueidis
func NewRueidis(cfg *conf.Bootstrap, logger log.Logger) rueidis.Client {
	l := log.NewHelper(log.With(logger, "module", "rueidis"))
	r, err := rueidiscache.NewRueidis(&rueidis.ClientOption{
		Username:    cfg.Data.Redis.Username,
		Password:    cfg.Data.Redis.Password,
		InitAddress: []string{cfg.Data.Redis.Addr},
		SelectDB:    int(cfg.Data.Redis.Db),
	})
	if err != nil {
		l.Fatalf("failed opening connection to redis")
		panic(fmt.Sprintf("NewRueidis err: %s", err))
	}
	return r
}
