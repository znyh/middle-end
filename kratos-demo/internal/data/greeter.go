package data

import (
	"context"
	"fmt"
	"time"

	"kratos-demo/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type greeterRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewGreeterRepo(data *Data, logger log.Logger) biz.GreeterRepo {
	return &greeterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *greeterRepo) Save(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {

	key := "users:"
	field := fmt.Sprintf("%+v", g.NickName)
	val := fmt.Sprintf("%+v", g.NickName)
	if err := r.data.redis.HExists(ctx, key, field).Err(); err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	if err := r.data.redis.HSet(ctx, key, field, val).Err(); err != nil {
		return nil, err
	}
	if err := r.data.redis.Expire(ctx, key, time.Hour*24*30).Err(); err != nil {
		return nil, err
	}

	var user biz.Greeter
	// 验证是否已经创建
	result := r.data.db.Where(&biz.Greeter{NickName: g.NickName}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	user.NickName = g.NickName
	if res := r.data.db.Create(&user); res.Error != nil {
		return nil, status.Errorf(codes.Internal, res.Error.Error())
	}

	return &user, nil
}

func (r *greeterRepo) Update(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	return g, nil
}

func (r *greeterRepo) FindByID(context.Context, int64) (*biz.Greeter, error) {
	return nil, nil
}

func (r *greeterRepo) ListByHello(context.Context, string) ([]*biz.Greeter, error) {
	return nil, nil
}

func (r *greeterRepo) ListAll(context.Context) ([]*biz.Greeter, error) {
	return nil, nil
}
