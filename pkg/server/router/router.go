package router

import (
	"context"

	"github.com/threefoldtech/grid3-go/deployer"
)

type Router struct {
	Routes map[string]func(r *Router, ctx context.Context, data string) (interface{}, error)
	Client *deployer.TFPluginClient
}
