//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/biz"
	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/conf"
	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/data"
	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/server"
	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.JWT, log.Logger) (*kratos.App, func(), *service.BPMService, error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		newApp,
	))
}
