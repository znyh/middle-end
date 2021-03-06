// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-cqrs/app/logger/service/internal/biz"
	"kratos-cqrs/app/logger/service/internal/conf"
	"kratos-cqrs/app/logger/service/internal/data"
	"kratos-cqrs/app/logger/service/internal/server"
	"kratos-cqrs/app/logger/service/internal/service"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(confServer *conf.Server, registry *conf.Registry, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	client := data.NewEntClient(confData, logger)
	dataData, cleanup, err := data.NewData(client, logger)
	if err != nil {
		return nil, nil, err
	}
	sensorDataRepo := data.NewSensorDataRepo(dataData, logger)
	sensorDataUseCase := biz.NewSensorDataUseCase(sensorDataRepo, logger)
	sensorRepo := data.NewSensorRepo(dataData, logger)
	sensorUseCase := biz.NewSensorUseCase(sensorRepo, logger)
	loggerService := service.NewLoggerService(sensorDataUseCase, sensorUseCase, logger)
	grpcServer := server.NewGRPCServer(confServer, logger, loggerService)
	registrar := server.NewConsulRegistrar(registry)
	app := newApp(logger, grpcServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}
