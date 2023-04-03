package router

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

type ZDB struct {
	NodeID      uint32 `json:"node_id"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Public      bool   `json:"public"`
	Size        int    `json:"size"`
	Description string `json:"description"`
	Mode        string `json:"mode"`
	Port        uint32 `json:"port"`
	Namespace   string `json:"namespace"`

	// computed
	IPs []string `json:"ips"`
}

func (r *Router) ZDBDeploy(ctx context.Context, data string) (interface{}, error) {
	zdb := ZDB{}

	if err := json.Unmarshal([]byte(data), &zdb); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal zdb model data")
	}

	originalProjectName := zdb.Name
	cliProjectName := generateProjectName(zdb.Name)
	zdb.Name = cliProjectName

	zdb, err := r.deployZDB(ctx, zdb)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deploy zdb")
	}

	zdb.Name = originalProjectName

	return zdb, nil
}

func (r *Router) ZDBDelete(ctx context.Context, data string) (interface{}, error) {
	var zdbName string

	if err := json.Unmarshal([]byte(data), &zdbName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal zdb model data")
	}

	cliProjectName := generateProjectName(zdbName)
	zdbName = cliProjectName

	err := r.deleteZDB(ctx, zdbName)
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

	originalProjectName := zdbName
	cliProjectName := generateProjectName(zdbName)
	zdbName = cliProjectName

	zdb, err := r.getZDB(ctx, zdbName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get zdb")
	}

	zdb.Name = originalProjectName

	return zdb, nil
}

///

func (r *Router) deployZDB(ctx context.Context, zdb ZDB) (ZDB, error) {
	// validate no workloads with the same name
	if err := r.validateProjectName(ctx, zdb.Name); err != nil {
		return ZDB{}, err
	}

	// deploy
	zdbs := []workloads.ZDB{
		NewClientWorkloadFromZDB(zdb),
	}
	log.Info().Msgf("Deploying zdb: %+v", zdbs)

	clientDeployment := workloads.NewDeployment(zdb.Name, zdb.NodeID, zdb.Name, nil, "", nil, zdbs, nil, nil)
	err := r.Client.DeploymentDeployer.Deploy(ctx, &clientDeployment)
	if err != nil {
		return ZDB{}, errors.Wrapf(err, "failed to deploy zdb with name: %s", zdb.Name)
	}

	// get the result with the computed values
	zdb_, err := r.Client.State.LoadZdbFromGrid(zdb.NodeID, zdb.Name, zdb.Name)
	result := NewZDBFromClientZDB(zdb_)
	result.NodeID = zdb.NodeID

	// NOTE: clean the state after deploying
	return result, nil
}

func (r *Router) deleteZDB(ctx context.Context, name string) error {
	// TODO: fix canceling
	err := r.Client.CancelByProjectName(name)
	if err != nil {
		errors.Wrapf(err, "Failed to cancel cluster with name: %s", name)
	}

	return nil
}

func (r *Router) getZDB(ctx context.Context, name string) (ZDB, error) {
	// get the contract
	contracts, err := r.Client.ContractsGetter.ListContractsOfProjectName(name)
	if err != nil {
		return ZDB{}, errors.Wrapf(err, "Couldn't get contract for name: %s", name)
	}

	result := ZDB{}
	for _, contract := range contracts.NodeContracts {

		cl, err := r.Client.NcPool.GetNodeClient(r.Client.SubstrateConn, contract.NodeID)
		if err != nil {
			return ZDB{}, errors.Wrapf(err, "Couldn't get client for node: %s", contract.NodeID)
		}

		cid, err := strconv.ParseUint(contract.ContractID, 10, 64)
		if err != nil {
			return ZDB{}, errors.Wrapf(err, "Couldn't parse contract Id: %s", contract.ContractID)
		}

		dl, err := cl.DeploymentGet(ctx, cid)
		if err != nil {
			return ZDB{}, errors.Wrapf(err, "Couldn't get deployment for contract Id: %s", contract.ContractID)
		}

		for _, workload := range dl.Workloads {
			if workload.Type == zos.ZDBType {
				zdb := workloads.ZDB{}

				zdb, err = workloads.NewZDBFromWorkload(&workload)
				if err != nil {
					return ZDB{}, errors.Wrapf(err, "Failed to get vm from workload: %s", workload)
				}

				result = NewZDBFromClientZDB(zdb)
				result.NodeID = contract.NodeID
			}
		}
	}
	return result, nil
}

func NewClientWorkloadFromZDB(zdb ZDB) workloads.ZDB {
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

func NewZDBFromClientZDB(wl workloads.ZDB) ZDB {
	return ZDB{
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
