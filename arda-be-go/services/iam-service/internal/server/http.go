package server

import (
	stdlib "net/http"

	"github.com/arda-labs/arda/pkg/middleware"
	"github.com/arda-labs/arda/services/iam-service/internal/conf"
	"github.com/arda-labs/arda/services/iam-service/internal/service"
	pb "github.com/arda-labs/arda/services/iam-service/api/iam/v1"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

func NewHTTPServer(c *conf.Server, iam *service.IAMService, menu *service.MenuService, logger log.Logger) *khttp.Server {
	var opts = []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
			middleware.Logging(logger),
			middleware.Auth(),
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
		out, err := menu.GetMenu(ctx, &service.GetMenuRequest{})
		if err != nil {
			return err
		}
		return ctx.Result(200, out)
	})
	srv.Route("/").GET("/v1/menus", func(ctx khttp.Context) error {
		in := &service.ListMenusRequest{}
		_ = ctx.BindQuery(in)
		out, err := menu.ListMenus(ctx, in)
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
		out, err := menu.CreateMenu(ctx, in)
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
		out, err := menu.UpdateMenu(ctx, in)
		if err != nil {
			return err
		}
		return ctx.Result(200, out)
	})
	srv.Route("/").DELETE("/v1/menus/{id}", func(ctx khttp.Context) error {
		in := &service.DeleteMenuRequest{}
		_ = ctx.BindVars(in)
		out, err := menu.DeleteMenu(ctx, in)
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
