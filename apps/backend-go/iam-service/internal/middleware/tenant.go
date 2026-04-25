package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

func TenantExtractor() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			query := tr.RequestHeader().Get("X-Tenant-ID")
			if query != "" && TenantIDFromContext(ctx) == "" {
				ctx = context.WithValue(ctx, tenantIDKey, query)
			}
			return handler(ctx, req)
		}
	}
}
