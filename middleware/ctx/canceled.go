package ctx

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

var RequestCanceledErr = errors.New(http.StatusConflict, "RequestCanceledErr", "request canceled")

func Canceled() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			panicChan := make(chan any, 1)
			done := make(chan struct{})
			go func() {
				defer func() {
					if p := recover(); p != nil {
						panicChan <- p
					}
				}()
				reply, err = handler(ctx, req)
				close(done)
			}()
			select {
			case p := <-panicChan:
				panic(p)
			case <-ctx.Done():
				if errors.Is(ctx.Err(), context.Canceled) {
					return nil, RequestCanceledErr
				}
				return nil, ctx.Err()
			case <-done:
				return reply, err
			}
		}
	}
}
