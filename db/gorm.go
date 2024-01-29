package db

import (
	"fmt"

	conf "github.com/fzf-labs/fkratos-contrib/api/conf/v1"

	"github.com/fzf-labs/fpkg/orm"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// NewGorm 初始化gorm
func NewGorm(cfg *conf.Bootstrap, logger log.Logger) *gorm.DB {
	l := log.NewHelper(log.With(logger, "module", "NewGorm"))
	client, err := orm.NewGormPostgresClient(&orm.GormPostgresClientConfig{
		DataSourceName:  cfg.Data.Gorm.DataSourceName,
		MaxIdleConn:     int(cfg.Data.Gorm.MaxIdleConn),
		MaxOpenConn:     int(cfg.Data.Gorm.MaxOpenConn),
		ConnMaxLifeTime: cfg.Data.Gorm.ConnMaxLifeTime.AsDuration(),
		ShowLog:         cfg.Data.Gorm.ShowLog,
		Tracing:         cfg.Data.Gorm.Tracing,
	})
	if err != nil {
		l.Fatalf("failed opening connection to postgres")
		panic(fmt.Sprintf("NewGorm postgres err: %s", err))
	}
	return client
}
