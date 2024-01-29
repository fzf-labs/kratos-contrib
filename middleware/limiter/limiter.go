package limiter

import (
	conf "github.com/fzf-labs/fkratos-contrib/api/conf/v1"

	"github.com/go-kratos/aegis/ratelimit/bbr"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
)

// Limit 限流器
func Limit(rateLimiter *conf.RateLimiter) middleware.Middleware {
	return ratelimit.Server(
		ratelimit.WithLimiter(
			bbr.NewLimiter([]bbr.Option{
				bbr.WithWindow(rateLimiter.Window.AsDuration()),
				bbr.WithBucket(int(rateLimiter.Bucket)),
				bbr.WithCPUThreshold(rateLimiter.CpuThreshold),
			}...),
		),
	)
}
