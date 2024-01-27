package bootstrap

import (
	"context"
	conf "fkratos-contrib/api/conf/v1"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/registry"
	kGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"

	"time"
)

const defaultTimeout = 5 * time.Second

// NewGrpcClient 创建GRPC客户端
func NewGrpcClient(ctx context.Context, r registry.Discovery, discovery, serviceName string, timeoutDuration *durationpb.Duration) *grpc.ClientConn {
	timeout := defaultTimeout
	if timeoutDuration != nil {
		timeout = timeoutDuration.AsDuration()
	}
	endpoint := "discovery:///"
	switch discovery {
	case "nacos":
		serviceName += ".grpc"
	case "polaris":
		serviceName += "grpc"
	default:
	}
	endpoint += serviceName
	conn, err := kGrpc.DialInsecure(
		ctx,
		kGrpc.WithEndpoint(endpoint),
		kGrpc.WithDiscovery(r),
		kGrpc.WithTimeout(timeout),
		kGrpc.WithMiddleware(
			recovery.Recovery(),
			metadata.Client(),
			tracing.Client(),
			circuitbreaker.Client(),
		),
	)
	if err != nil {
		log.Fatalf("dial grpc client [%s] failed: %s", serviceName, err.Error())
	}
	return conn
}

// NewGrpcServer 创建GRPC服务端
func NewGrpcServer(cfg *conf.Bootstrap, logger log.Logger, m ...middleware.Middleware) *kGrpc.Server {
	var opts []kGrpc.ServerOption
	var ms []middleware.Middleware
	ms = append(ms,
		tracing.Server(),
		logging.Server(logger),
		recovery.Recovery(),
		metadata.Server(),
		validate.Validator(),
		Metrics(),
	)
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
	return kGrpc.NewServer(opts...)
}
