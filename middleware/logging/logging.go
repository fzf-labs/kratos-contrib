package logging

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// Redacter defines how to log an object
type Redacter interface {
	Redact() string
}

// Server is an server logging middleware.
func Server(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			level, stack := extractError(err)
			_ = log.WithContext(ctx, logger).Log(level,
				"kind", "server",
				"component", kind,
				"operation", operation,
				"req", extractReq(req),
				"reply", extractReply(reply),
				"code", code,
				"reason", reason,
				"stack", stack,
				"latency", time.Since(startTime).Milliseconds(),
			)
			return
		}
	}
}

// Client is a client logging middleware.
func Client(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromClientContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			level, stack := extractError(err)
			_ = log.WithContext(ctx, logger).Log(level,
				"kind", "client",
				"component", kind,
				"operation", operation,
				"req", extractReq(req),
				"reply", extractReply(reply),
				"code", code,
				"reason", reason,
				"stack", stack,
				"latency", time.Since(startTime).Milliseconds(),
			)
			return
		}
	}
}

// extractArgs returns the string of the req
func extractReq(req interface{}) string {
	if redactor, ok := req.(Redacter); ok {
		return redactor.Redact()
	}
	m, _ := sonic.MarshalString(req)
	if sonic.Valid([]byte(m)) {
		return m
	}
	return fmt.Sprintf("%+v", req)
}

// extractReply returns the string of the reply
func extractReply(reply interface{}) string {
	m, _ := sonic.MarshalString(reply)
	return m
}

// extractError returns the string of the error
func extractError(err error) (level log.Level, str string) {
	if err != nil {
		return log.LevelError, fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, ""
}
