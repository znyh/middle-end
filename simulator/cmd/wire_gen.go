// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"simulator/internal/biz"
	"simulator/internal/conf"
	"simulator/internal/data"
	"simulator/internal/server"
	"simulator/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	client := data.NewRedis(confData)
	dataData, cleanup, err := data.NewData(confData, client, logger)
	if err != nil {
		return nil, nil, err
	}
	useRepo := data.NewUseRepo(dataData, logger)
	usecase := biz.NewUsecase(useRepo, logger)
	simulatorService := service.NewSimulatorService(usecase)
	grpcServer := server.NewGRPCServer(confServer, simulatorService, logger)
	httpServer := server.NewHTTPServer(confServer, simulatorService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
