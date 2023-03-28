package router

import "context"

func GatewayNameDeploy(ctx context.Context, data string) (string, error)

func GatewayNameGet(ctx context.Context, data string) (string, error)

func GatewayNameDelete(ctx context.Context, data string) (string, error)

func GatewayFQDNDeploy(ctx context.Context, data string) (string, error)

func GatewayFQDNGet(ctx context.Context, data string) (string, error)

func GatewayFQDNDelete(ctx context.Context, data string) (string, error)
