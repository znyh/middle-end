package service

import (
	"context"

	pb "calculator/api/calculator/v1"
	"calculator/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type CalculatorService struct {
	pb.UnimplementedCalculatorServer

	uc  *biz.CalculatorUseCase
	log *log.Helper
}

func NewCalculatorService(uc *biz.CalculatorUseCase, logger log.Logger) *CalculatorService {
	return &CalculatorService{uc: uc, log: log.NewHelper(logger)}
}

func (s *CalculatorService) Calc(ctx context.Context, req *pb.CalcReq) (*pb.CalcRsp, error) {
	return s.uc.Calc(ctx, req)
}
