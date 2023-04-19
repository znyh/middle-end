package comet

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/go-kratos/kratos/pkg/ecode"
	"github.com/go-kratos/kratos/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// recovery is a server interceptor that recovers from any panics.
func (s *Server) recovery() UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) ([]byte, error) {
		var (
			err  error
			resp []byte
		)
		defer func() {
			if rerr := recover(); rerr != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				rs := runtime.Stack(buf, false)
				if rs > size {
					rs = size
				}
				buf = buf[:rs]
				pl := fmt.Sprintf("comet server panic: %v\n%v\n%s\n", req, rerr, buf)
				fmt.Fprintf(os.Stderr, pl)
				log.Error(pl)
				err = status.Errorf(codes.Unknown, ecode.ServerErr.Error())
			}
		}()
		resp, err = handler(ctx, req)
		return resp, err
	}
}
