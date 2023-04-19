package comet

import (
	"context"
)

type UnaryServerInfo struct {
	// Server is the service implementation the user provides. This is read-only.
	Server interface{}
	// FullMethod is the full RPC method string, i.e., /package.service/method.
	FullMethod string
}

type UnaryHandler func(ctx context.Context, req interface{}) ([]byte, error)

type UnaryServerInterceptor func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) ([]byte, error)
