package router

import (
	"context"
)

type Router struct {
	Routes map[string]func(r *Router, ctx context.Context, data string) (interface{}, error)
	// Client *deployer.TFPluginClient
	client TFGridClient
}
