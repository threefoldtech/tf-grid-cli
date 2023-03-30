package router

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/deployer"
	procedure "github.com/threefoldtech/tf-grid-cli/pkg/server/procedures"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
)

func GatewayNameDeploy(ctx context.Context, client *deployer.TFPluginClient, data string) (interface{}, error) {
	model := types.GatewayNameModel{}
	if err := json.Unmarshal([]byte(data), &model); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model data")
	}

	res, err := procedure.GatewayNameDeploy(ctx, model, client)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to deploy gateway %s", model.Name)
	}

	return res, nil
}

func GatewayNameGet(ctx context.Context, client *deployer.TFPluginClient, data string) (interface{}, error) {
	modelName := ""
	if err := json.Unmarshal([]byte(data), &modelName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model name")
	}

	res, err := procedure.GatewayNameGet(ctx, modelName, client)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to delete gateway model %s", modelName)
	}

	return res, nil
}

func GatewayNameDelete(ctx context.Context, client *deployer.TFPluginClient, data string) (interface{}, error) {
	modelName := ""
	if err := json.Unmarshal([]byte(data), &modelName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model name")
	}

	if err := procedure.GatewayNameDelete(ctx, modelName, client); err != nil {
		return nil, errors.Wrapf(err, "failed to delete gateway model %s", modelName)
	}

	return nil, nil
}

func GatewayFQDNDeploy(ctx context.Context, client *deployer.TFPluginClient, data string) (interface{}, error) {
	model := types.GatewayFQDNModel{}
	if err := json.Unmarshal([]byte(data), &model); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model data")
	}

	res, err := procedure.GatewayFQDNDeploy(ctx, model, client)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to deploy gateway %s", model.Name)
	}

	return res, nil
}

func GatewayFQDNGet(ctx context.Context, client *deployer.TFPluginClient, data string) (interface{}, error) {
	modelName := ""
	if err := json.Unmarshal([]byte(data), &modelName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model name")
	}

	res, err := procedure.GatewayFQDNGet(ctx, modelName, client)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to delete gateway model %s", modelName)
	}

	return res, nil
}

func GatewayFQDNDelete(ctx context.Context, client *deployer.TFPluginClient, data string) (interface{}, error) {
	modelName := ""
	if err := json.Unmarshal([]byte(data), &modelName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model name")
	}

	if err := procedure.GatewayFQDNDelete(ctx, modelName, client); err != nil {
		return nil, errors.Wrapf(err, "failed to delete gateway model %s", modelName)
	}

	return nil, nil
}
