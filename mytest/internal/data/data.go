package data

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "mytest/internal/conf"
    "mytest/pkg/kafka"

    klog "github.com/go-kratos/kratos/v2/log"
    "github.com/go-redis/redis/extra/redisotel"
    "github.com/go-redis/redis/v8"
    "github.com/google/wire"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "gorm.io/gorm/schema"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewKafka, NewUseRepo)

// Data .
type Data struct {
    db    *gorm.DB
    redis *redis.Client
    pub   kafka.Producer
}

// NewData .
func NewData(c *conf.Data, db *gorm.DB, r *redis.Client, pub kafka.Producer) (*Data, func(), error) {
    d := &Data{db: db, redis: r, pub: pub}
    cleanup := func() {
        klog.Info("closing the data resources")
        if _, err := d.db.DB(); err != nil {
            klog.Error(err)
        }
        if err := d.redis.Close(); err != nil {
            klog.Error(err)
        }
    }
    return d, cleanup, nil
}

// NewDB .
func NewDB(c *conf.Data) *gorm.DB {
    // 终端打印输入 sql 执行记录
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
        logger.Config{
            SlowThreshold: time.Second, // 慢查询 SQL 阈值
            Colorful:      true,        // 禁用彩色打印
            //IgnoreRecordNotFoundError: false,
            LogLevel: logger.Info, // Log lever
        },
    )

    db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
        Logger:                                   newLogger,
        DisableForeignKeyConstraintWhenMigrating: true,
        NamingStrategy: schema.NamingStrategy{
            //SingularTable: true, // 表名是否加 s
        },
    })

    if err != nil {
        klog.Errorf("failed opening connection to sqlite: %v", err)
        panic("failed to connect database")
    }

    return db
}

// NewRedis .
func NewRedis(c *conf.Data) *redis.Client {
    r := redis.NewClient(&redis.Options{
        Addr:         c.Redis.Addr,
        Password:     c.Redis.Password,
        DB:           int(c.Redis.Db),
        DialTimeout:  c.Redis.DialTimeout.AsDuration(),
        WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
        ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
    })
    r.AddHook(redisotel.TracingHook{})
    if err := r.Ping(context.Background()).Err(); err != nil {
        klog.Errorf("failed opening connection to redis: %v", err)
        panic("failed to connect redis")
    }
    return r
}

// NewKafka .
func NewKafka(c *conf.Data) (producer kafka.Producer, err error) {
    p, err := kafka.NewProducer(c.Kafka.Endpoints)
    if err != nil {
        klog.Errorf("failed opening connection to kafka(%+v) err:%+v", c.Kafka.Endpoints, err)
        panic("failed to connect kafka")
    }

    pingKafka(p, c)

    return p, nil
}

func pingKafka(producer kafka.Producer, c *conf.Data) {

    //produce message
    if err := producer.Producer(kafka.Message{
        Topic: "test",
        Value: []byte("wo shi yi tiao yu"),
    }); err != nil {
        klog.Errorf("failed to start Producer: %v . ", err)
        panic(err)
    }
    klog.Infof("start Producer...")

    //consumer message
    go func() {
        var (
            handler = func(msg []byte, args ...interface{}) {
                if len(args) != 1 {
                    panic("arg is empty,")
                }
                value, ok := args[0].(int)
                if !ok {
                    panic("args 1 not type int")
                }
                klog.Infof("kafka consumer. value:%+v msg:%+v", value, string(msg))
            }
        )

        consumer, err := kafka.NewConsumer(c.Kafka.Endpoints, "zhuma")
        if err != nil {
            klog.Errorf("failed to start consumer, err:", err)
            panic("failed to start consumer")
        }
        defer consumer.Close()

        klog.Infof("start consumer...")
        if err = consumer.Consume(map[string]kafka.Handler{
            "test": kafka.Handler{
                Run:  handler,
                Args: []interface{}{1},
            },
        }); err != nil {
            panic(fmt.Sprintf("kafka consumer err:%+v", err))
        }
    }()

    return
}
