package biz

import (
    "context"

    v1 "mytest/api"
)

type (
    UseRepo interface {
        Save(ctx context.Context) error
    }

    UseCase struct {
        repo UseRepo
    }
)

func NewUsecase(repo UseRepo) *UseCase {
    return &UseCase{repo: repo}
}

func (u *UseCase) SayHello(ctx context.Context, req *v1.HelloRequest) (*v1.HelloReply, error) {
    u.repo.Save(ctx)
    return &v1.HelloReply{}, nil
}
