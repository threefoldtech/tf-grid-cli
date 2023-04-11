package router

import (
	"context"

	client "github.com/threefoldtech/tf-grid-cli/pkg/server/cli_client"
)

type Router struct {
	Routes map[string]func(r *Router, ctx context.Context, data string) (interface{}, error)
	// Client *deployer.TFPluginClient
	client client.CLIClient
}
