package biz

import (
    "context"

    v1 "simulator/api/simulator/v1"

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

func (uc *Usecase) OnBetReq(ctx context.Context, in *v1.BetReq) (*v1.BetRsp, error) {
    rsp := &v1.BetRsp{
        GameID: in.GameID,
        Uid:    in.Uid,
        Data:   in.Data,
    }
    uc.log.WithContext(ctx).Infof("OnBetReq: %v rsp: %+v", in, rsp)
    return rsp, nil
}
