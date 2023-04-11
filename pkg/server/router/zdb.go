package router

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	client "github.com/threefoldtech/tf-grid-cli/pkg/server/cli_client"
)

func (r *Router) ZDBDeploy(ctx context.Context, data string) (interface{}, error) {
	zdb := client.ZDB{}

	if err := json.Unmarshal([]byte(data), &zdb); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal zdb model data")
	}

	projectName := generateProjectName(zdb.Name)

	zdb, err := r.client.ZDBDeploy(ctx, zdb, projectName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deploy zdb")
	}

	return zdb, nil
}

func (r *Router) ZDBDelete(ctx context.Context, data string) (interface{}, error) {
	var zdbName string

	if err := json.Unmarshal([]byte(data), &zdbName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal zdb model data")
	}

	projectName := generateProjectName(zdbName)

	err := r.client.ZDBDelete(ctx, projectName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to delete zdb")
	}

	return nil, nil
}

func (r *Router) ZDBGet(ctx context.Context, data string) (interface{}, error) {
	var zdbName string

	if err := json.Unmarshal([]byte(data), &zdbName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal zdb model data")
	}

	projectName := generateProjectName(zdbName)

	zdb, err := r.client.ZDBGet(ctx, projectName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get zdb")
	}

	return zdb, nil
}
