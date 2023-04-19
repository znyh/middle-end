package registrar

import (
	"kratos-demo/internal/conf"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
	etcdclient "go.etcd.io/etcd/client/v3"
)

// NewEtcdRegistrar 引入 etcd
func NewEtcdRegistrar(conf *conf.Registry) registry.Registrar {
	endPoints := conf.Etcd.Endpoints
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: endPoints,
	})
	if err != nil {
		panic(err)
	}
	r := etcd.New(client)

	return r
}
