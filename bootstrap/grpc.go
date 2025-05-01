package bootstrap

import (
	"context"
	"time"

	conf "github.com/fzf-labs/kratos-contrib/api/conf/v1"
	"github.com/fzf-labs/kratos-contrib/middleware/limiter"
	"github.com/fzf-labs/kratos-contrib/middleware/logging"
	"github.com/fzf-labs/kratos-contrib/middleware/metrics"
	"github.com/go-kratos/kratos/contrib/middleware/validate/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	kGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
)

const defaultTimeout = 5 * time.Second

// NewGrpcClient 创建GRPC客户端
func NewGrpcClient(
	ctx context.Context,
	cfg *conf.Bootstrap,
	logger log.Logger,
	serverName string,
	r registry.Discovery,
	m ...middleware.Middleware,
) *grpc.ClientConn {
	timeout := defaultTimeout
	endpoint := "discovery:///" + serverName
	var ms []middleware.Middleware
	if cfg.GetClient() != nil && cfg.GetClient().GetGrpc() != nil {
		if cfg.GetClient().GetGrpc().GetTimeout() != nil {
			timeout = cfg.GetClient().GetGrpc().GetTimeout().AsDuration()
		}
		if cfg.GetClient().GetGrpc().GetMiddleware() != nil {
			if cfg.GetClient().GetGrpc().GetMiddleware().GetEnableTracing() {
				ms = append(ms, tracing.Client())
			}
			if cfg.GetClient().GetGrpc().GetMiddleware().GetEnableLogging() {
				ms = append(ms, logging.Client(logger))
			}
			if cfg.GetClient().GetGrpc().GetMiddleware().GetEnableRecovery() {
				ms = append(ms, recovery.Recovery())
			}
			if cfg.GetClient().GetGrpc().GetMiddleware().GetEnableCircuitBreaker() {
				ms = append(ms, circuitbreaker.Client())
			}
			if cfg.GetClient().GetGrpc().GetMiddleware().GetEnableMetadata() {
				ms = append(ms, metadata.Client())
			}
			if cfg.GetClient().GetGrpc().GetMiddleware().GetEnableValidate() {
				ms = append(ms, validate.ProtoValidate())
			}
		}
	}
	ms = append(ms, m...)
	conn, err := kGrpc.DialInsecure(
		ctx,
		kGrpc.WithEndpoint(endpoint),
		kGrpc.WithDiscovery(r),
		kGrpc.WithTimeout(timeout),
		kGrpc.WithMiddleware(ms...),
	)
	if err != nil {
		log.Fatalf("dial grpc client [%s] failed: %s", serverName, err.Error())
	}
	return conn
}

// NewGrpcServer 创建GRPC服务端
func NewGrpcServer(
	cfg *conf.Bootstrap,
	logger log.Logger,
	m ...middleware.Middleware,
) *kGrpc.Server {
	var opts []kGrpc.ServerOption
	var ms []middleware.Middleware
	if cfg.GetServer() != nil && cfg.GetServer().GetGrpc() != nil && cfg.GetServer().GetGrpc().GetMiddleware() != nil {
		if cfg.GetServer().GetGrpc().GetMiddleware().GetEnableTracing() {
			ms = append(ms, tracing.Server())
		}
		if cfg.GetServer().GetGrpc().GetMiddleware().GetEnableLogging() {
			ms = append(ms, logging.Server(logger))
		}
		if cfg.GetServer().GetGrpc().GetMiddleware().GetEnableRecovery() {
			ms = append(ms, recovery.Recovery())
		}
		if cfg.GetServer().GetGrpc().GetMiddleware().GetEnableMetrics() {
			ms = append(ms, metrics.Server())
		}
		if cfg.GetServer().GetGrpc().GetMiddleware().GetEnableRateLimiter() {
			ms = append(ms, limiter.Limit(cfg.GetServer().GetGrpc().GetMiddleware().GetLimiter()))
		}
		if cfg.GetServer().GetGrpc().GetMiddleware().GetEnableMetadata() {
			ms = append(ms, metadata.Server())
		}
		if cfg.GetServer().GetGrpc().GetMiddleware().GetEnableValidate() {
			ms = append(ms, validate.ProtoValidate())
		}
	}
	ms = append(ms, m...)
	opts = append(opts, kGrpc.Middleware(ms...))
	if cfg.GetServer().GetGrpc().GetNetwork() != "" {
		opts = append(opts, kGrpc.Network(cfg.GetServer().GetGrpc().GetNetwork()))
	}
	if cfg.GetServer().GetGrpc().GetAddr() != "" {
		opts = append(opts, kGrpc.Address(cfg.GetServer().GetGrpc().GetAddr()))
	}
	if cfg.GetServer().GetGrpc().GetTimeout() != nil {
		opts = append(opts, kGrpc.Timeout(cfg.GetServer().GetGrpc().GetTimeout().AsDuration()))
	}
	srv := kGrpc.NewServer(opts...)
	return srv
}
