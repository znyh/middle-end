package service

import (
    "context"

    v1 "mytest/api"

    "mytest/internal/biz"
)

type Service struct {
    v1.UnimplementedGreeterServer
    uc *biz.UseCase
}

func NewGreeterService(uc *biz.UseCase) *Service {
    return &Service{uc: uc}
}

func (s *Service) SayHello(ctx context.Context, req *v1.HelloRequest) (*v1.HelloReply, error) {
    return s.uc.SayHello(ctx, req)
}
