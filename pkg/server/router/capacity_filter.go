package router

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/deployer"
	procedure "github.com/threefoldtech/tf-grid-cli/pkg/server/procedures"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
)

func FilterNodes(ctx context.Context, client *deployer.TFPluginClient, data string) (interface{}, error) {
	options := types.FilterOptions{}

	err := json.Unmarshal([]byte(data), &options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal filter options")
	}

	res, err := procedure.FilterNodes(ctx, options, client)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to filter nodes for options %s", options)
	}

	return res, nil
}
