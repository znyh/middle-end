package comet

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/pkg/ecode"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/go-kratos/kratos/pkg/net/metadata"
)

// Comet Log Flag
const (
	// disable all log.
	LogFlagDisable = 1 << iota
	// disable print args on log.
	LogFlagDisableArgs
	// disable info level log.
	LogFlagDisableInfo
)

func logFn(code int, dt time.Duration) func(context.Context, ...log.D) {
	switch {
	case code < 0:
		return log.Errorv
	case dt >= time.Millisecond*500:
		// TODO: slowlog make it configurable.
		return log.Warnv
	case code > 0:
		return log.Warnv
	}
	return log.Infov
}

// serverLogging warden grpc logging
func serverLogging(logFlag int8) UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) ([]byte, error) {
		startTime := time.Now()
		caller := metadata.String(ctx, metadata.Caller)
		if caller == "" {
			caller = "no_user"
		}
		// call server handler
		resp, err := handler(ctx, req)

		// after server response
		code := ecode.Cause(err).Code()
		duration := time.Since(startTime)
		// monitor
		_metricServerReqDur.Observe(int64(duration/time.Millisecond), info.FullMethod, caller)
		_metricServerReqCodeTotal.Inc(info.FullMethod, caller, strconv.Itoa(code))

		if logFlag&LogFlagDisable != 0 {
			return resp, err
		}
		// TODO: find better way to deal with slow log.
		if logFlag&LogFlagDisableInfo != 0 && err == nil && duration < 500*time.Millisecond {
			return resp, err
		}
		logFields := []log.D{
			log.KVString("user", caller),
			log.KVString("path", info.FullMethod),
			log.KVInt("ret", code),
			log.KVFloat64("ts", duration.Seconds()),
			log.KVString("source", "comet-access-log"),
		}
		if logFlag&LogFlagDisableArgs == 0 {
			// TODO: it will panic if someone remove String method from protobuf message struct that auto generate from protoc.
			logFields = append(logFields, log.KVString("args", req.(fmt.Stringer).String()))
		}
		if code < 0 {
			logFields = append(logFields, log.KVString("error", err.Error()), log.KVString("stack", fmt.Sprintf("%+v", err)))
			logFn(code, duration)(ctx, logFields...)
		}
		return resp, err
	}
}
