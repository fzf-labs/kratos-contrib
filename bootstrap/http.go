package bootstrap

import (
	conf "fkratos-contrib/api/conf/v1"
	"fkratos-contrib/middleware/limiter"
	"fkratos-contrib/middleware/logging"
	"net/http/pprof"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
)

// NewHTTPServer 创建Http服务端
func NewHTTPServer(cfg *conf.Bootstrap, logger log.Logger, m ...middleware.Middleware) *http.Server {
	var opts []http.ServerOption
	var ms []middleware.Middleware
	if cfg.Server != nil && cfg.Server.Http != nil && cfg.Server.Http.Middleware != nil {
		if cfg.Server.Http.Middleware.GetEnableTracing() {
			ms = append(ms, tracing.Server())
		}
		if cfg.Server.Http.Middleware.GetEnableRecovery() {
			ms = append(ms, recovery.Recovery())
		}
		if cfg.Server.Http.Middleware.GetEnableLogging() {
			ms = append(ms, logging.Server(logger))
		}
		if cfg.Client.Http.Middleware.GetEnableMetadata() {
			ms = append(ms, metadata.Client())
		}
		if cfg.Server.Http.Middleware.GetEnableRateLimiter() {
			ms = append(ms, limiter.Limit(cfg.Server.Http.Middleware.Limiter))
		}
		if cfg.Server.Http.Middleware.GetEnableValidate() {
			ms = append(ms, validate.Validator())
		}
	}
	ms = append(ms, m...)
	opts = append(opts, http.Middleware(ms...))
	if cfg.Server.Http.GetEnableCors() {
		opts = append(opts, http.Filter(handlers.CORS(
			handlers.AllowedHeaders(cfg.Server.Http.Cors.Headers),
			handlers.AllowedMethods(cfg.Server.Http.Cors.Methods),
			handlers.AllowedOrigins(cfg.Server.Http.Cors.Origins),
		)))
	}
	if cfg.Server.Http.Network != "" {
		opts = append(opts, http.Network(cfg.Server.Http.Network))
	}
	if cfg.Server.Http.Addr != "" {
		opts = append(opts, http.Address(cfg.Server.Http.Addr))
	}
	if cfg.Server.Http.Timeout != nil {
		opts = append(opts, http.Timeout(cfg.Server.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	if cfg.Server.Http.GetEnablePprof() {
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
