package server

import (
	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"

	"github.com/yygqzzk/review-service/internal/conf"
)

// 注册服务
func NewRegistry(conf *conf.Registry) (registry.Registrar, error) {
	config := api.DefaultConfig()
	config.Address = conf.Consul.Address
	config.Scheme = conf.Consul.Scheme

	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	return consul.New(client, consul.WithHealthCheck(true)), nil
}

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, NewRegistry)
