package biz

import (
    "context"

    "github.com/go-kratos/kratos/v2/log"
)

// UseRepo is a Use repo.
type UseRepo interface {
    Save(context.Context, ) (err error)
}

// Usecase is a usecase.
type Usecase struct {
    repo UseRepo
    log  *log.Helper
}

// NewUsecase new a usecase.
func NewUsecase(repo UseRepo, logger log.Logger) *Usecase {
    return &Usecase{repo: repo, log: log.NewHelper(logger)}
}

// SayHello .
func (uc *Usecase) SayHello(ctx context.Context) error {
    uc.log.WithContext(ctx).Infof("CreateGreeter: %v", "g.Hello")
    return uc.repo.Save(ctx)
}
