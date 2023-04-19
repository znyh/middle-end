package data

import (
	"calculator/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type calculatorRepo struct {
	data *Data
	log  *log.Helper
}

// NewCalculatorRepo .
func NewCalculatorRepo(data *Data, logger log.Logger) biz.CalculatorRepo {
	return &calculatorRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
