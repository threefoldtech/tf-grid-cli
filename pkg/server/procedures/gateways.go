package procedure

import (
	"context"

	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
)

func GatewayNameDeploy(ctx context.Context, gatewayNameModel types.GatewayNameModel) (types.GatewayNameModel, error)

func GatewayNameDelete(ctx context.Context, name string) error

func GatewayNameGet(ctx context.Context, name string) (types.GatewayNameModel, error)

func GatewayFQDNDeploy(ctx context.Context, gatewayFQDNModel types.GatewayFQDNModel) (types.GatewayFQDNModel, error)

func GatewayFQDNDelete(ctx context.Context, name string) error

func GatewayFQDNGet(ctx context.Context, name string) (types.GatewayFQDNModel, error)
