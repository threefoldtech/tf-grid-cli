package router

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/deployer"
	procedure "github.com/threefoldtech/tf-grid-cli/pkg/server/procedures"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
)

func ZDBDeploy(ctx context.Context, client *deployer.TFPluginClient, data string) (interface{}, error) {
	zdb := types.ZDB{}

	if err := json.Unmarshal([]byte(data), &zdb); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal zdb model data")
	}

	zdb, err := procedure.ZDBDeploy(ctx, zdb, client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deploy zdb")
	}

	return zdb, nil
}

func ZDBDelete(ctx context.Context, client *deployer.TFPluginClient, data string) (interface{}, error) {
	var zdbName string

	if err := json.Unmarshal([]byte(data), &zdbName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal zdb model data")
	}

	err := procedure.ZDBDelete(ctx, zdbName, client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to delete zdb")
	}

	return nil, nil
}

func ZDBGet(ctx context.Context, client *deployer.TFPluginClient, data string) (interface{}, error) {
	var zdbName string

	if err := json.Unmarshal([]byte(data), &zdbName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal zdb model data")
	}

	zdb, err := procedure.ZDBGet(ctx, zdbName, client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get zdb")
	}

	return zdb, nil
}
