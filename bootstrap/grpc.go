package bootstrap

import (
	"context"
	"time"

	conf "github.com/fzf-labs/fkratos-contrib/api/conf/v1"
	"github.com/fzf-labs/fkratos-contrib/middleware/limiter"
	"github.com/fzf-labs/fkratos-contrib/middleware/logging"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
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
	if cfg.Client != nil && cfg.Client.Grpc != nil {
		if cfg.Client.Grpc.Timeout != nil {
			timeout = cfg.Client.Grpc.Timeout.AsDuration()
		}
		if cfg.Client.Grpc.Middleware != nil {
			if cfg.Client.Grpc.Middleware.GetEnableTracing() {
				ms = append(ms, tracing.Client())
			}
			if cfg.Client.Grpc.Middleware.GetEnableRecovery() {
				ms = append(ms, recovery.Recovery())
			}
			if cfg.Client.Grpc.Middleware.GetEnableLogging() {
				ms = append(ms, logging.Client(logger))
			}
			if cfg.Client.Grpc.Middleware.GetEnableMetadata() {
				ms = append(ms, metadata.Client())
			}
			if cfg.Client.Grpc.Middleware.GetEnableCircuitBreaker() {
				ms = append(ms, circuitbreaker.Client())
			}
			if cfg.Client.Grpc.Middleware.GetEnableValidate() {
				ms = append(ms, validate.Validator())
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
	if cfg.Server != nil && cfg.Server.Grpc != nil && cfg.Server.Grpc.Middleware != nil {
		if cfg.Server.Grpc.Middleware.GetEnableTracing() {
			ms = append(ms, tracing.Server())
		}
		if cfg.Server.Grpc.Middleware.GetEnableRecovery() {
			ms = append(ms, recovery.Recovery())
		}
		if cfg.Server.Grpc.Middleware.GetEnableLogging() {
			ms = append(ms, logging.Server(logger))
		}
		if cfg.Client.Grpc.Middleware.GetEnableMetadata() {
			ms = append(ms, metadata.Client())
		}
		if cfg.Server.Grpc.Middleware.GetEnableRateLimiter() {
			ms = append(ms, limiter.Limit(cfg.Server.Grpc.Middleware.Limiter))
		}
		if cfg.Server.Grpc.Middleware.GetEnableValidate() {
			ms = append(ms, validate.Validator())
		}
	}
	ms = append(ms, m...)
	opts = append(opts, kGrpc.Middleware(ms...))
	if cfg.Server.Grpc.Network != "" {
		opts = append(opts, kGrpc.Network(cfg.Server.Grpc.Network))
	}
	if cfg.Server.Grpc.Addr != "" {
		opts = append(opts, kGrpc.Address(cfg.Server.Grpc.Addr))
	}
	if cfg.Server.Grpc.Timeout != nil {
		opts = append(opts, kGrpc.Timeout(cfg.Server.Grpc.Timeout.AsDuration()))
	}
	srv := kGrpc.NewServer(opts...)
	return srv
}
