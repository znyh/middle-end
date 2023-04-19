package data

import (
    "context"

    "mytest/internal/biz"
)

type useRepo struct {
    data *Data
}

// NewUseRepo .
func NewUseRepo(data *Data) biz.UseRepo {
    return &useRepo{
        data: data,
    }
}

func (r *useRepo) Save(context.Context) error {
    return nil
}
