package service

import (
    "context"

    v1 "simulator/api/simulator/v1"
    "simulator/internal/biz"
)

// SimulatorService is a Simulator service.
type SimulatorService struct {
    v1.UnimplementedSimulatorServer

    uc *biz.Usecase
}

// NewSimulatorService new a greeter service.
func NewSimulatorService(uc *biz.Usecase) *SimulatorService {
    return &SimulatorService{uc: uc}
}

// SayHello implements simulator.SimulatorServer.
func (s *SimulatorService) SayHello(ctx context.Context, in *v1.HelloReq) (*v1.HelloRsp, error) {
    err := s.uc.SayHello(ctx)
    if err != nil {
        return nil, err
    }
    return &v1.HelloRsp{Message: "Hello " + in.Name}, nil
}

// OnBetReq implements OnBetReq
func (s *SimulatorService) OnBetReq(ctx context.Context, in *v1.BetReq) (*v1.BetRsp, error) {
    rsp, err := s.uc.OnBetReq(ctx, in)
    return rsp, err
}

// OnCancelBetReq implements OnCancelBetReq
func (s *SimulatorService) OnCancelBetReq(context.Context, *v1.CancelBetReq) (*v1.CancelBetRsp, error) {
    return nil, nil
}

// OnGetBetListReq implements OnGetBetListReq
func (s *SimulatorService) OnGetBetListReq(context.Context, *v1.GetBetListReq) (*v1.GetBetListRsp, error) {
    return nil, nil
}
