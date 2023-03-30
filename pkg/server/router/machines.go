package router

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/grid3-go/deployer"
	procedure "github.com/threefoldtech/tf-grid-cli/pkg/server/procedures"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
)

func MachinesDeploy(ctx context.Context, client *deployer.TFPluginClient, data string) (interface{}, error) {
	model := types.MachinesModel{}

	if err := json.Unmarshal([]byte(data), &model); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal machine model data")
	}

	model, err := procedure.MachinesDeploy(ctx, model, client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deploy model")
	}

	return model, nil
}

func MachinesDelete(ctx context.Context, client *deployer.TFPluginClient, data string) (interface{}, error) {
	modelName := ""
	err := json.Unmarshal([]byte(data), &modelName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal model name")
	}

	log.Info().Msgf("cancelilng project %s", modelName)
	if err := procedure.MachinesDelete(ctx, modelName, client); err != nil {
		return nil, errors.Wrapf(err, "failed to delete model %s", modelName)
	}

	return nil, nil
}

func MachinesGet(ctx context.Context, client *deployer.TFPluginClient, data string) (interface{}, error) {
	modelName := ""
	err := json.Unmarshal([]byte(data), &modelName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal model name")
	}

	log.Info().Msgf("getting project %s", modelName)
	model, err := procedure.MachinesGet(ctx, modelName, client)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get model %s", modelName)
	}

	return model, nil
}

// func MachineAdd(ctx context.Context, data string) (string, error)

// func MachineRemove(ctx context.Context, data string) (string, error)
