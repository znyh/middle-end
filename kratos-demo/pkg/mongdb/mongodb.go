//package mongdb
//
//import (
//	"time"
//
//	"git.huoys.com/middle-end/kratos/pkg/net/netutil/breaker"
//	"github.com/go-kratos/kratos/v2/log"
//)
//
//// Config sqlserver config.
//type Config struct {
//    URI            string          // likes mongodb://foo:bar@localhost:27017
//    ConnectTimeout time.Duration   // connection mongodb timeout
//    QueryTimeout   time.Duration   // query mongodb timeout
//    ExecTimeout    time.Duration   // execute mongodb timeout
//    Breaker        *breaker.Config // breaker
//}
//
//// NewMongoDB new db and retry connection when has error.
//func NewMongoDB(c *Config) (db *DB) {
//    if c.ConnectTimeout == 0 || c.QueryTimeout == 0 || c.ExecTimeout == 0 {
//        panic("mongo must be set query/execute/connect timeout")
//    }
//    db, err := Open(c)
//    if err != nil {
//        log.Error("open mongodb error(%v)", err)
//        panic(err)
//    }
//
//    return
//}

package mongdb
