//go:build wireinject
// +build wireinject

package main

import (
	"iam-service/internal/biz"
	"iam-service/internal/conf"
	"iam-service/internal/data"
	"iam-service/internal/server"
	"iam-service/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func wireApp(*conf.Server, *conf.Data, *conf.JWT, *conf.Zitadel, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
