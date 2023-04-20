package data

import (
    "context"

    "simulator/internal/biz"

    "github.com/go-kratos/kratos/v2/log"
)

type useRepo struct {
    data *Data
    log  *log.Helper
}

// NewUseRepo .
func NewUseRepo(data *Data, logger log.Logger) biz.UseRepo {
    return &useRepo{
        data: data,
        log:  log.NewHelper(logger),
    }
}

func (r *useRepo) Save(ctx context.Context) (err error) {
    return nil
}
