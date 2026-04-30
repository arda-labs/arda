package server

import (
	"context"
	stdlib "net/http"
	"strings"

	"github.com/arda-labs/arda/arda-be-go/pkg/middleware"
	notificationv1 "github.com/arda-labs/arda/arda-be-go/services/notification-service/api/notification/v1"
	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/conf"
	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/service"
	kratoserrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	kmiddleware "github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

func NewHTTPServer(c *conf.Server, jwt *conf.JWT, svc *service.NotificationService, logger log.Logger) *khttp.Server {
	var opts = []khttp.ServerOption{
		khttp.Filter(apiPrefixAlias()),
		khttp.Middleware(
			recovery.Recovery(),
			middleware.Logging(logger),
			middleware.Auth(middleware.WithJWTValidation(jwt.JwksEndpoint, jwt.Issuer, jwt.Audience)),
			authRequired(),
			middleware.Tenant(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, khttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, khttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, khttp.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := khttp.NewServer(opts...)
	notificationv1.RegisterNotificationServiceHTTPServer(srv, svc)
	srv.HandleFunc("/health", func(w stdlib.ResponseWriter, r *stdlib.Request) {
		w.WriteHeader(stdlib.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	return srv
}

func apiPrefixAlias() khttp.FilterFunc {
	return func(next stdlib.Handler) stdlib.Handler {
		return stdlib.HandlerFunc(func(w stdlib.ResponseWriter, r *stdlib.Request) {
			if r.URL != nil && (r.URL.Path == "/api/v1" || strings.HasPrefix(r.URL.Path, "/api/v1/")) {
				clone := r.Clone(r.Context())
				urlCopy := *r.URL
				urlCopy.Path = strings.TrimPrefix(r.URL.Path, "/api")
				if urlCopy.RawPath != "" {
					urlCopy.RawPath = strings.TrimPrefix(urlCopy.RawPath, "/api")
				}
				clone.URL = &urlCopy
				next.ServeHTTP(w, clone)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func authRequired() kmiddleware.Middleware {
	return func(handler kmiddleware.Handler) kmiddleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if r, ok := khttp.RequestFromServerContext(ctx); ok && r.URL != nil && r.URL.Path == "/health" {
				return handler(ctx, req)
			}
			if middleware.GetUserID(ctx) == "" {
				return nil, kratoserrors.Unauthorized("UNAUTHORIZED", "missing subject")
			}
			return handler(ctx, req)
		}
	}
}
