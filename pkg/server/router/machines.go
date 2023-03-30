package router

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	procedure "github.com/threefoldtech/tf-grid-cli/pkg/server/procedures"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
)

func MachinesDeploy(ctx context.Context, data string) (interface{}, error) {
	model := types.MachinesModel{}

	if err := json.Unmarshal([]byte(data), &model); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal machine model data")
	}

	client, err := getClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new grid client")
	}

	model, err = procedure.MachinesDeploy(ctx, model, &client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deploy model")
	}

	return model, nil
}

func MachinesDelete(ctx context.Context, data string) (interface{}, error) {
	modelName := ""
	err := json.Unmarshal([]byte(data), &modelName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal model name")
	}

	client, err := getClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new grid client")
	}

	log.Info().Msgf("cancelilng project %s", modelName)
	if err := procedure.MachinesDelete(ctx, modelName, &client); err != nil {
		return nil, errors.Wrapf(err, "failed to delete model %s", modelName)
	}

	return struct{}{}, err
}

// func MachinesGet(ctx context.Context, data string) ([]byte, error) {

// }

// func MachineAdd(ctx context.Context, data string) (string, error)

// func MachineRemove(ctx context.Context, data string) (string, error)
