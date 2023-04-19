package main

import (
    "context"
    "fmt"
    "log"

    pb "github.com/go-kratos/examples/helloworld/helloworld"
    "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
    "github.com/go-kratos/kratos/v2"
    "github.com/go-kratos/kratos/v2/middleware/recovery"
    "github.com/go-kratos/kratos/v2/transport/grpc"
    "github.com/go-kratos/kratos/v2/transport/http"
    etcdclient "go.etcd.io/etcd/client/v3"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
    pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
    return &pb.HelloReply{Message: fmt.Sprintf("Welcome %+v!", in.Name)}, nil
}

func main() {
    client, err := etcdclient.New(etcdclient.Config{
        Endpoints: []string{"127.0.0.1:2379"},
    })
    if err != nil {
        log.Fatal(err)
    }

    httpSrv := http.NewServer(
        http.Address(":8010"),
        http.Middleware(
            recovery.Recovery(),
        ),
    )
    grpcSrv := grpc.NewServer(
        grpc.Address(":9010"),
        grpc.Middleware(
            recovery.Recovery(),
        ),
    )

    s := &server{}
    pb.RegisterGreeterServer(grpcSrv, s)
    pb.RegisterGreeterHTTPServer(httpSrv, s)

    r := etcd.New(client)
    app := kratos.New(
        kratos.Name("helloworld"),
        kratos.Server(
            httpSrv,
            grpcSrv,
        ),
        kratos.Registrar(r),
    )
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
