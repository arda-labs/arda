package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

const (
	userIDKey          contextKey = "user_id"
	userEmailKey       contextKey = "user_email"
	userPermissionsKey contextKey = "user_permissions"
)

// Auth middleware extracts user information from headers (set by APISIX/Zitadel)
func Auth() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				userID := tr.RequestHeader().Get("X-User-Id")
				if userID != "" {
					ctx = context.WithValue(ctx, userIDKey, userID)
					ctx = context.WithValue(ctx, userEmailKey, tr.RequestHeader().Get("X-User-Email"))
					ctx = context.WithValue(ctx, userPermissionsKey, tr.RequestHeader().Get("X-User-Permissions"))
				}
			}
			return handler(ctx, req)
		}
	}
}

func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(userIDKey).(string); ok {
		return v
	}
	return ""
}

func GetEmail(ctx context.Context) string {
	if v, ok := ctx.Value(userEmailKey).(string); ok {
		return v
	}
	return ""
}
