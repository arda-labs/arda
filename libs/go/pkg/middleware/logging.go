package middleware

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// Logging middleware logs request details with context info
func Logging(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if tr, ok := transport.FromServerContext(ctx); ok {
				kind = tr.Kind().String()
				operation = tr.Operation()
			}
			reply, err := handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			level := log.LevelInfo
			if err != nil {
				level = log.LevelError
			}

			_ = log.WithContext(ctx, logger).Log(level,
				"kind", kind,
				"operation", operation,
				"args", req,
				"code", code,
				"reason", reason,
				"latency", time.Since(startTime).String(),
				"user_id", GetUserID(ctx),
				"tenant_id", GetTenantID(ctx),
			)
			return reply, err
		}
	}
}
