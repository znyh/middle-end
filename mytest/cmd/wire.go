//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
    "mytest/internal/biz"
    "mytest/internal/conf"
    "mytest/internal/data"
    "mytest/internal/server"
    "mytest/internal/service"

    "github.com/google/wire"
)

func wireApp(*conf.Server, *conf.Data) (*App, func(), error) {
    panic(wire.Build(biz.ProviderSet, data.ProviderSet, service.ProviderSet, server.ProviderSet, newApp))
}
