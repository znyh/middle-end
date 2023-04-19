package server

import (
    v1 "mytest/api"
    "mytest/internal/conf"
    "mytest/internal/service"

    "google.golang.org/grpc"
)

func NewGrpcServer(c *conf.Server, svc *service.Service) *grpc.Server {

    // 创建一个服务
    srv := grpc.NewServer()
    // 注册服务
    v1.RegisterGreeterServer(srv, svc)

    return srv
}
