package data

import (
    "context"

    "simulator/internal/conf"

    "github.com/go-kratos/kratos/v2/log"
    "github.com/go-redis/redis/extra/redisotel"
    "github.com/go-redis/redis/v8"
    "github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUseRepo, NewRedis)

// Data .
type Data struct {
    // TODO wrapped database client
    redis *redis.Client
}

// NewData .
func NewData(c *conf.Data, r *redis.Client, logger log.Logger) (*Data, func(), error) {
    d := &Data{redis: r}
    cleanup := func() {
        log.NewHelper(logger).Info("closing the data resources")
        if err := d.redis.Close(); err != nil {
            log.Error(err)
        }
    }

    return d, cleanup, nil
}

// NewRedis .
func NewRedis(c *conf.Data) *redis.Client {
    r := redis.NewClient(&redis.Options{
        Addr: c.Redis.Addr,
        //Password:     c.Redis.Password,
        //DB:           int(c.Redis.Db),
        //DialTimeout:  c.Redis.DialTimeout.AsDuration(),
        WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
        ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
    })
    r.AddHook(redisotel.TracingHook{})
    if err := r.Ping(context.Background()).Err(); err != nil {
        log.Errorf("failed opening connection to redis: %v", err)
        panic("failed to connect redis")
    }
    return r
}
