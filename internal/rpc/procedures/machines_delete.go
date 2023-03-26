package procedure

import (
	"context"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/deployer"
)

func MachinesDelete(ctx context.Context, projectName string, client deployer.TFPluginClient) error {
	err := client.CancelByProjectName(projectName)
	if err != nil {
		return errors.Wrap(err, "failed to cancel project")
	}

	return nil
}
