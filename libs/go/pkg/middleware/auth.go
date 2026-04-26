package middleware

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v5"
)

const (
	userIDKey          contextKey = "user_id"
	userEmailKey       contextKey = "user_email"
	userPermissionsKey contextKey = "user_permissions"
)

// Auth middleware extracts user information from headers (set by APISIX/Zitadel)
// or parses Bearer token directly for local development.
func Auth() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				userID := tr.RequestHeader().Get("X-User-Id")
				email := tr.RequestHeader().Get("X-User-Email")
				perms := tr.RequestHeader().Get("X-User-Permissions")

				// Fallback for local development: Parse Bearer token
				if userID == "" {
					authHeader := tr.RequestHeader().Get("Authorization")
					if strings.HasPrefix(authHeader, "Bearer ") {
						tokenStr := authHeader[7:]
						// Parse unverified JWT for local dev
						// In production, verification happens at APISIX/Zitadel
						token, _, _ := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
						if token != nil {
							if claims, ok := token.Claims.(jwt.MapClaims); ok {
								// Extract standard claims
								if sub, ok := claims["sub"].(string); ok {
									userID = sub
								}
								if mail, ok := claims["email"].(string); ok {
									email = mail
								}
								// Extract permissions from Zitadel custom claims
								if p, ok := claims["urn:zitadel:iam:org:project:roles"].(map[string]any); ok {
									var permList []string
									for role := range p {
										permList = append(permList, role)
									}
									perms = strings.Join(permList, ",")
								} else if p, ok := claims["roles"].([]any); ok {
									var permList []string
									for _, r := range p {
										if role, ok := r.(string); ok {
											permList = append(permList, role)
										}
									}
									perms = strings.Join(permList, ",")
								}
								// Extract tenant ID if available
								if tenantID := tr.RequestHeader().Get("X-Tenant-ID"); tenantID == "" {
									if tid, ok := claims["tenant_id"].(string); ok {
										tr.RequestHeader().Set("X-Tenant-ID", tid)
									} else if tid, ok := claims["urn:zitadel:iam:org:id"].(string); ok {
										tr.RequestHeader().Set("X-Tenant-ID", tid)
									}
								}
							}
						}
					}
				}

				if userID != "" {
					ctx = context.WithValue(ctx, userIDKey, userID)
					ctx = context.WithValue(ctx, userEmailKey, email)
					ctx = context.WithValue(ctx, userPermissionsKey, perms)
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
