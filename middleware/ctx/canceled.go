package ctx

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	RequestCanceledErr = errors.New(http.StatusConflict, "RequestCanceledErr", "request canceled")
	RequestTimeoutErr  = errors.New(http.StatusConflict, "RequestTimeoutErr", "request timeout")
)

// Canceled 用于处理请求取消或超时的情况
func Canceled(timeout time.Duration) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			st := time.Now()
			reply, err = handler(ctx, req)
			if ctxErr := ctx.Err(); ctxErr != nil {
				if errors.Is(ctxErr, context.Canceled) || (errors.Is(ctxErr, context.DeadlineExceeded) && time.Since(st) < timeout) {
					return nil, RequestCanceledErr
				}
				if errors.Is(ctxErr, context.DeadlineExceeded) {
					return nil, RequestTimeoutErr
				}
			}
			if err != nil {
				grpcErr, ok := status.FromError(err)
				if ok {
					switch grpcErr.Code() {
					case codes.Canceled:
						return nil, RequestCanceledErr
					case codes.DeadlineExceeded:
						return nil, RequestTimeoutErr
					}
				}
			}
			return reply, err
		}
	}
}
