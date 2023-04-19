package mssql

import (
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/go-kratos/kratos/pkg/net/netutil/breaker"
	"github.com/go-kratos/kratos/pkg/time"

	// database driver
	_ "github.com/denisenkom/go-mssqldb"
)

// Config sqlserver config.
type Config struct {
	Host              string   // write data source name.
	ReadHost          []string // read data source name.
	User              string
	Password          string
	ProcPassword      string
	Database          string
	EncodeConfig      bool
	EncryptEnabled    bool
	LogFlags          int             //logging flags (default 0/no logging, 63 for full logging)
	ConnectionTimeout int             // in seconds
	Active            int             // pool
	Idle              int             // pool
	IdleTimeout       time.Duration   // connect max life time.
	QueryTimeout      time.Duration   // query sql timeout
	ExecTimeout       time.Duration   // execute sql timeout
	TranTimeout       time.Duration   // transaction sql timeout
	Breaker           *breaker.Config // breaker
}

// NewSQLServer new db and retry connection when has error.
func NewSQLServer(c *Config) (db *DB) {
	if c.QueryTimeout == 0 || c.ExecTimeout == 0 || c.TranTimeout == 0 {
		panic("sqlserver must be set query/execute/transction timeout")
	}
	db, err := Open(c)
	if err != nil {
		log.Error("open sqlserver error(%v)", err)
		panic(err)
	}
	return
}
