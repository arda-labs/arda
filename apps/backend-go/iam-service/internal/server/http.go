package server

import (
	"context"
	stdlib "net/http"
	"strings"

	"github.com/arda-labs/arda/arda-be-go/pkg/middleware"
	pb "github.com/arda-labs/arda/arda-be-go/services/iam-service/api/iam/v1"
	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/conf"
	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
)

func NewHTTPServer(c *conf.Server, jwt *conf.JWT, iam *service.IAMService, menu *service.MenuService, logger log.Logger) *khttp.Server {
	var opts = []khttp.ServerOption{
		khttp.Filter(
			apiPrefixAlias(),
			handlers.CORS(
				handlers.AllowedOrigins([]string{
					"http://localhost:3000", "http://127.0.0.1:3000",
					"http://localhost:3001", "http://127.0.0.1:3001",
					"http://localhost:3002", "http://127.0.0.1:3002",
					"http://localhost:4200", "http://127.0.0.1:4200",
				}),
				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
				handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Tenant-ID", "X-Request-ID", "Accept-Language", "X-User-Username"}),
			),
		),
		khttp.Middleware(
			recovery.Recovery(),
			middleware.Logging(logger),
			middleware.Auth(middleware.WithJWTValidation(jwt.JwksEndpoint, jwt.Issuer, jwt.Audience)),
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
	pb.RegisterIAMServiceHTTPServer(srv, iam)

	// Menu routes (plain JSON, no proto codegen needed)
	srv.Route("/").GET("/v1/me/menu", func(ctx khttp.Context) error {
		in := &service.GetMenuRequest{}
		khttp.SetOperation(ctx, "/iam.v1.MenuService/GetMenu")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return menu.GetMenu(ctx, req.(*service.GetMenuRequest))
		})
		out, err := h(ctx, in)
		if err != nil {
			return err
		}
		return ctx.Result(200, out)
	})
	srv.Route("/").GET("/v1/menus", func(ctx khttp.Context) error {
		in := &service.ListMenusRequest{}
		_ = ctx.BindQuery(in)
		khttp.SetOperation(ctx, "/iam.v1.MenuService/ListMenus")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return menu.ListMenus(ctx, req.(*service.ListMenusRequest))
		})
		out, err := h(ctx, in)
		if err != nil {
			return err
		}
		return ctx.Result(200, out)
	})
	srv.Route("/").GET("/v1/users/{user_id}/groups", func(ctx khttp.Context) error {
		in := &service.ListUserGroupsRequest{}
		_ = ctx.BindQuery(in)
		in.UserID = ctx.Vars()["user_id"][0]
		khttp.SetOperation(ctx, "/iam.v1.IAMService/ListUserGroups")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return iam.ListUserGroups(ctx, req.(*service.ListUserGroupsRequest))
		})
		out, err := h(ctx, in)
		if err != nil {
			return err
		}
		return ctx.Result(200, out)
	})
	srv.Route("/").GET("/v1/users/{user_id}/effective-permissions", func(ctx khttp.Context) error {
		in := &service.GetUserEffectivePermissionsRequest{}
		_ = ctx.BindQuery(in)
		in.UserID = ctx.Vars()["user_id"][0]
		khttp.SetOperation(ctx, "/iam.v1.IAMService/GetUserEffectivePermissions")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return iam.GetUserEffectivePermissions(ctx, req.(*service.GetUserEffectivePermissionsRequest))
		})
		out, err := h(ctx, in)
		if err != nil {
			return err
		}
		return ctx.Result(200, out)
	})
	srv.Route("/").POST("/v1/menus", func(ctx khttp.Context) error {
		in := &service.CreateMenuRequest{}
		if err := ctx.Bind(in); err != nil {
			return err
		}
		khttp.SetOperation(ctx, "/iam.v1.MenuService/CreateMenu")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return menu.CreateMenu(ctx, req.(*service.CreateMenuRequest))
		})
		out, err := h(ctx, in)
		if err != nil {
			return err
		}
		return ctx.Result(200, out)
	})
	srv.Route("/").PUT("/v1/menus/{id}", func(ctx khttp.Context) error {
		in := &service.UpdateMenuRequest{}
		if err := ctx.Bind(in); err != nil {
			return err
		}
		_ = ctx.BindVars(in)
		khttp.SetOperation(ctx, "/iam.v1.MenuService/UpdateMenu")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return menu.UpdateMenu(ctx, req.(*service.UpdateMenuRequest))
		})
		out, err := h(ctx, in)
		if err != nil {
			return err
		}
		return ctx.Result(200, out)
	})
	srv.Route("/").DELETE("/v1/menus/{id}", func(ctx khttp.Context) error {
		in := &service.DeleteMenuRequest{}
		_ = ctx.BindVars(in)
		khttp.SetOperation(ctx, "/iam.v1.MenuService/DeleteMenu")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return menu.DeleteMenu(ctx, req.(*service.DeleteMenuRequest))
		})
		out, err := h(ctx, in)
		if err != nil {
			return err
		}
		return ctx.Result(200, out)
	})

	// Health endpoint
	srv.HandleFunc("/health", func(w stdlib.ResponseWriter, r *stdlib.Request) {
		w.WriteHeader(stdlib.StatusOK)
		w.Write([]byte("OK"))
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
