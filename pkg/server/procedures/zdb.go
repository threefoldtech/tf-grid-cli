package procedure

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

func ZDBDeploy(ctx context.Context, zdb types.ZDB, client *deployer.TFPluginClient) (types.ZDB, error) {
	// validate no workloads with the same name
	contracts, err := client.ContractsGetter.ListContractsOfProjectName(zdb.Name)
	if err != nil {
		return types.ZDB{}, errors.Wrapf(err, "failed to retrieve contracts with project name %s", zdb.Name)
	}

	if len(contracts.NameContracts) > 0 || len(contracts.NodeContracts) > 0 || len(contracts.RentContracts) > 0 {
		return types.ZDB{}, fmt.Errorf("there is a zdb with the same name: %s", zdb.Name)
	}

	// deploy
	zdbs := []workloads.ZDB{
		convertZDBtoWorkload(zdb),
	}
	log.Info().Msgf("Deploying zdb: %+v", zdbs)

	clientDeployment := workloads.NewDeployment(zdb.Name, zdb.NodeID, zdb.Name, nil, "", nil, zdbs, nil, nil)
	err = client.DeploymentDeployer.Deploy(ctx, &clientDeployment)
	if err != nil {
		return types.ZDB{}, errors.Wrapf(err, "failed to deploy zdb with name: %s", zdb.Name)
	}

	// get the result with the computed values
	zdb_, err := client.State.LoadZdbFromGrid(zdb.NodeID, zdb.Name, zdb.Name)
	result := convertWrokloadtoZDB(zdb_)
	result.NodeID = zdb.NodeID

	// NOTE: clean the state after deploying
	return result, nil
}

func ZDBDelete(ctx context.Context, name string, client *deployer.TFPluginClient) error {
	// TODO: fix canceling
	err := client.CancelByProjectName(name)
	if err != nil {
		errors.Wrapf(err, "Failed to cancel cluster with name: %s", name)
	}

	return nil
}

func ZDBGet(ctx context.Context, name string, client *deployer.TFPluginClient) (types.ZDB, error) {
	// get the contract
	contracts, err := client.ContractsGetter.ListContractsOfProjectName(name)
	if err != nil {
		return types.ZDB{}, errors.Wrapf(err, "Couldn't get contract for name: %s", name)
	}

	result := types.ZDB{}
	for _, contract := range contracts.NodeContracts {

		cl, err := client.NcPool.GetNodeClient(client.SubstrateConn, contract.NodeID)
		if err != nil {
			return types.ZDB{}, errors.Wrapf(err, "Couldn't get client for node: %s", contract.NodeID)
		}

		cid, err := strconv.ParseUint(contract.ContractID, 10, 64)
		if err != nil {
			return types.ZDB{}, errors.Wrapf(err, "Couldn't parse contract Id: %s", contract.ContractID)
		}

		dl, err := cl.DeploymentGet(ctx, cid)
		if err != nil {
			return types.ZDB{}, errors.Wrapf(err, "Couldn't get deployment for contract Id: %s", contract.ContractID)
		}

		log.Info().Msgf("wl: %+v", dl.Workloads)

		for _, workload := range dl.Workloads {
			if workload.Type == zos.ZDBType {
				zdb := workloads.ZDB{}

				zdb, err = workloads.NewZDBFromWorkload(&workload)
				if err != nil {
					return types.ZDB{}, errors.Wrapf(err, "Failed to get vm from workload: %s", workload)
				}

				result = convertWrokloadtoZDB(zdb)
				result.NodeID = contract.NodeID
			}
		}
	}
	return result, nil
}

func convertZDBtoWorkload(zdb types.ZDB) workloads.ZDB {
	return workloads.ZDB{
		Name:        zdb.Name,
		Password:    zdb.Password,
		Public:      zdb.Public,
		Size:        zdb.Size,
		Description: zdb.Description,
		Mode:        zdb.Mode,
		Port:        zdb.Port,
		Namespace:   zdb.Namespace,
	}
}

func convertWrokloadtoZDB(wl workloads.ZDB) types.ZDB {
	return types.ZDB{
		Name:        wl.Name,
		Password:    wl.Password,
		Public:      wl.Public,
		Size:        wl.Size,
		Description: wl.Description,
		Mode:        wl.Mode,
		Port:        wl.Port,
		Namespace:   wl.Namespace,
		IPs:         wl.IPs,
	}
}
