package bootstrap

import (
	"net/http/pprof"

	conf "github.com/fzf-labs/kratos-contrib/api/conf/v1"
	"github.com/fzf-labs/kratos-contrib/middleware/limiter"
	"github.com/fzf-labs/kratos-contrib/middleware/logging"
	"github.com/fzf-labs/kratos-contrib/middleware/metrics"
	"github.com/go-kratos/kratos/contrib/middleware/validate/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewHTTPServer 创建Http服务端
func NewHTTPServer(cfg *conf.Bootstrap, logger log.Logger, m ...middleware.Middleware) *http.Server {
	var opts []http.ServerOption
	var ms []middleware.Middleware
	if cfg.GetServer() != nil && cfg.GetServer().GetHttp() != nil && cfg.GetServer().GetHttp().GetMiddleware() != nil {
		if cfg.GetServer().GetHttp().GetMiddleware().GetEnableTracing() {
			ms = append(ms, tracing.Server())
		}
		if cfg.GetServer().GetHttp().GetMiddleware().GetEnableLogging() {
			ms = append(ms, logging.Server(logger))
		}
		if cfg.GetServer().GetHttp().GetMiddleware().GetEnableRecovery() {
			ms = append(ms, recovery.Recovery())
		}
		if cfg.GetServer().GetHttp().GetMiddleware().GetEnableMetrics() {
			ms = append(ms, metrics.Server())
		}
		if cfg.GetServer().GetHttp().GetMiddleware().GetEnableRateLimiter() {
			ms = append(ms, limiter.Limit(cfg.GetServer().GetHttp().GetMiddleware().GetLimiter()))
		}
		if cfg.GetServer().GetHttp().GetMiddleware().GetEnableMetadata() {
			ms = append(ms, metadata.Server())
		}
		if cfg.GetServer().GetHttp().GetMiddleware().GetEnableValidate() {
			ms = append(ms, validate.ProtoValidate())
		}
	}
	ms = append(ms, m...)
	opts = append(opts, http.Middleware(ms...))
	if cfg.GetServer().GetHttp().GetEnableCors() {
		opts = append(opts, http.Filter(handlers.CORS(
			handlers.AllowedHeaders(cfg.GetServer().GetHttp().GetCors().GetHeaders()),
			handlers.AllowedMethods(cfg.GetServer().GetHttp().GetCors().GetMethods()),
			handlers.AllowedOrigins(cfg.GetServer().GetHttp().GetCors().GetOrigins()),
		)))
	}
	if cfg.GetServer().GetHttp().GetNetwork() != "" {
		opts = append(opts, http.Network(cfg.GetServer().GetHttp().GetNetwork()))
	}
	if cfg.GetServer().GetHttp().GetAddr() != "" {
		opts = append(opts, http.Address(cfg.GetServer().GetHttp().GetAddr()))
	}
	if cfg.GetServer().GetHttp().GetTimeout() != nil {
		opts = append(opts, http.Timeout(cfg.GetServer().GetHttp().GetTimeout().AsDuration()))
	}
	srv := http.NewServer(opts...)
	if cfg.GetServer().GetHttp().GetMiddleware().GetEnableMetrics() {
		registerHttpMetrics(srv)
	}
	if cfg.GetServer().GetHttp().GetEnablePprof() {
		registerHttpPprof(srv)
	}
	return srv
}

// registerHttpPprof 注册http pprof
func registerHttpPprof(s *http.Server) {
	s.HandleFunc("/debug/pprof", pprof.Index)
	s.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.HandleFunc("/debug/pprof/trace", pprof.Trace)
	s.HandleFunc("/debug/pprof/allocs", pprof.Handler("allocs").ServeHTTP)
	s.HandleFunc("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
	s.HandleFunc("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	s.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	s.HandleFunc("/debug/pprof/mutex", pprof.Handler("mutex").ServeHTTP)
	s.HandleFunc("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
}

// registerHttpMetrics 注册http metrics
func registerHttpMetrics(s *http.Server) {
	s.Handle("/metrics", promhttp.Handler())
}
