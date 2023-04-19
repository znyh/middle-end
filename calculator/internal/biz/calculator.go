package biz

import (
	"context"
	"fmt"
	"time"

	pb "calculator/api/calculator/v1"
	"calculator/internal/mc/colorcards/calc"

	"github.com/go-kratos/kratos/v2/log"
)

// CalculatorRepo is a Calculator repo.
type CalculatorRepo interface {
}

// CalculatorUseCase is a Calculator use case.
type CalculatorUseCase struct {
	repo CalculatorRepo
	log  *log.Helper
}

// NewCalculatorUseCase new a Calculator  use case.
func NewCalculatorUseCase(repo CalculatorRepo, logger log.Logger) *CalculatorUseCase {
	return &CalculatorUseCase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *CalculatorUseCase) Calc(ctx context.Context, req *pb.CalcReq) (*pb.CalcRsp, error) {
	cards := []int32{
		11, 12, 13, 14, 15, 16, 17, 18, 19,
		11, 12, 13, 14, 15, 16, 17, 18, 19,
		11, 12, 13, 14, 15, 16, 17, 18, 19,
		11, 12, 13, 14, 15, 16, 17, 18, 19,

		21, 22, 23, 24, 25, 26, 27, 28, 29,
		21, 22, 23, 24, 25, 26, 27, 28, 29,
		21, 22, 23, 24, 25, 26, 27, 28, 29,
		21, 22, 23, 24, 25, 26, 27, 28, 29,

		31, 32, 33, 34, 35, 36, 37, 38, 39,
		31, 32, 33, 34, 35, 36, 37, 38, 39,
		31, 32, 33, 34, 35, 36, 37, 38, 39,
		31, 32, 33, 34, 35, 36, 37, 38, 39,

		41, 42, 43, 44, 45, 46, 47,
		41, 42, 43, 44, 45, 46, 47,
		41, 42, 43, 44, 45, 46, 47,
		41, 42, 43, 44, 45, 46, 47,
	}

	start := time.Now()
	str := calc.Test(req, cards)
	use := time.Since(start).Milliseconds()
	log.Infof(" use:%+v/ms %+v", use, str)

	rsp := &pb.CalcRsp{}
	if len(str) > 0 {
		rsp.Conf = fmt.Sprintf("use:%+v/ms", use) + " " + str[0]
	}
	if len(str) > 1 {
		rsp.Desc1 = str[1]
	}
	if len(str) > 2 {
		rsp.Desc2 = str[2]
	}
	if len(str) > 3 {
		rsp.Desc3 = str[3]
	}
	if len(str) > 4 {
		rsp.Desc4 = str[4]
	}
	return rsp, nil
}
