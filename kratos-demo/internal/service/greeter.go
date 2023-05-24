package service

import (
    "context"
    "fmt"
    "time"

    v1 "kratos-demo/api/helloworld/v1"
    "kratos-demo/internal/biz"

    "github.com/go-kratos/kratos/v2/log"
)

// GreeterService is a greeter service.
type GreeterService struct {
    v1.UnimplementedGreeterServer

    uc *biz.GreeterUsecase

    log *log.Helper
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase, logger log.Logger) *GreeterService {
    return &GreeterService{uc: uc, log: log.NewHelper(logger)}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
    g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{NickName: in.Name})
    if err != nil {
        s.log.WithContext(ctx).Infof("SayHello error in:%+v err:%+v", in, err)
        return nil, err
    }
    return &v1.HelloReply{Message: "Hello " + g.NickName}, nil
}

// OnBetReq implements helloworld.GreeterServer. //OnBetReq(context.Context, *BetReq) (*BetRsp, error)
func (s *GreeterService) OnBetReq(ctx context.Context, in *v1.BetReq) (*v1.BetRsp, error) {
    s.log.WithContext(ctx).Infof("OnBetReq in:%+v ", in)
    resp := v1.BetRsp{
        GameID: in.GameID + 100000,
        Uid:    in.Uid + 1000,
        Data:   fmt.Sprintf("%+v - %+v", time.Now(), in.Data),
    }
    return &resp, nil
}
