package procedure

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

func MachinesGet(ctx context.Context, projectName string, client deployer.TFPluginClient) ([]gridtypes.Deployment, error) {
	contracts, err := client.ContractsGetter.ListContractsOfProjectName(projectName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list project contracts")
	}

	deployments := []gridtypes.Deployment{}
	for _, id := range contracts.NodeContracts {
		nodeClient, err := client.NcPool.GetNodeClient(client.SubstrateConn, id.NodeID)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get node %d client", id.NodeID)
		}
		contractID, err := strconv.Atoi(id.ContractID)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse contract id (%s)", id.ContractID)
		}

		dl, err := nodeClient.DeploymentGet(ctx, uint64(contractID))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get deployment with id %d", contractID)
		}
		deployments = append(deployments, dl)
	}

	return deployments, nil
}
