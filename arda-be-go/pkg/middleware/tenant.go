package middleware

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type contextKey string

const (
	tenantIDKey contextKey = "tenant_id"
)

// Tenant middleware extracts tenant ID from headers or metadata
func Tenant() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				tenantID := tr.RequestHeader().Get("X-Tenant-ID")
				if tenantID == "" {
					// Fallback to gRPC metadata or other sources if needed
				}
				if tenantID != "" {
					ctx = context.WithValue(ctx, tenantIDKey, tenantID)
				}
			}
			return handler(ctx, req)
		}
	}
}

// GetTenantID returns the tenant ID from context
func GetTenantID(ctx context.Context) string {
	if tenantID, ok := ctx.Value(tenantIDKey).(string); ok {
		return tenantID
	}
	return ""
}
