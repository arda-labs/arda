//go:build wireinject
// +build wireinject

package main

import (
	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/biz"
	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/conf"
	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/data"
	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/server"
	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func wireApp(*conf.Server, *conf.Data, *conf.JWT, *conf.Zitadel, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
