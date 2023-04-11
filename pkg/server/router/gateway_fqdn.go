package router

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	client "github.com/threefoldtech/tf-grid-cli/pkg/server/cli_client"
)

func (r *Router) GatewayFQDNDeploy(ctx context.Context, data string) (interface{}, error) {
	model := client.GatewayFQDNModel{}
	if err := json.Unmarshal([]byte(data), &model); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model data")
	}

	projectName := generateProjectName(model.Name)

	res, err := r.client.GatewayFQDNDeploy(ctx, model, projectName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to deploy gateway %s", model.Name)
	}

	return res, nil
}

func (r *Router) GatewayFQDNGet(ctx context.Context, data string) (interface{}, error) {
	modelName := ""
	if err := json.Unmarshal([]byte(data), &modelName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model name")
	}

	projectName := generateProjectName(modelName)

	res, err := r.client.GatewayFQDNGet(ctx, projectName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to delete gateway model %s", modelName)
	}

	return res, nil
}

func (r *Router) GatewayFQDNDelete(ctx context.Context, data string) (interface{}, error) {
	modelName := ""
	if err := json.Unmarshal([]byte(data), &modelName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model name")
	}

	projectName := generateProjectName(modelName)

	if err := r.client.GatewayNameDelete(ctx, projectName); err != nil {
		return nil, errors.Wrapf(err, "failed to delete gateway model %s", modelName)
	}

	return nil, nil
}
