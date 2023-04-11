package router

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	client "github.com/threefoldtech/tf-grid-cli/pkg/server/cli_client"
)

func (r *Router) MachinesDeploy(ctx context.Context, data string) (interface{}, error) {
	model := client.MachinesModel{}

	if err := json.Unmarshal([]byte(data), &model); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal machine model data")
	}

	projectName := generateProjectName(model.Name)

	model, err := r.client.MachinesDeploy(ctx, model, projectName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deploy model")
	}

	return model, nil
}

func (r *Router) MachinesDelete(ctx context.Context, data string) (interface{}, error) {
	modelName := ""
	err := json.Unmarshal([]byte(data), &modelName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal model name")
	}

	projectName := generateProjectName(modelName)

	log.Info().Msgf("cancelilng project %s", projectName)
	if err := r.client.MachinesDelete(ctx, projectName); err != nil {
		return nil, errors.Wrapf(err, "failed to delete model %s", modelName)
	}

	return nil, nil
}

func (r *Router) MachinesGet(ctx context.Context, data string) (interface{}, error) {
	modelName := ""
	err := json.Unmarshal([]byte(data), &modelName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal model name")
	}

	projectName := generateProjectName(modelName)

	log.Info().Msgf("getting project %s", projectName)
	model, err := r.client.MachinesGet(ctx, modelName, projectName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get model %s", modelName)
	}

	return model, nil
}
