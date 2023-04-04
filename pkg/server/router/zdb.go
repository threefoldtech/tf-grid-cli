package router

import (
	"context"
	"encoding/json"
	"fmt"
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

	// computed
	Port      uint32   `json:"port"`
	Namespace string   `json:"namespace"`
	IPs       []string `json:"ips"`
}

func (r *Router) ZDBDeploy(ctx context.Context, data string) (interface{}, error) {
	zdb := ZDB{}

	if err := json.Unmarshal([]byte(data), &zdb); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal zdb model data")
	}

	projectName := generateProjectName(zdb.Name)

	zdb, err := r.deployZDB(ctx, zdb, projectName)
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

	err := r.deleteZDB(ctx, projectName)
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

	zdb, err := r.getZDB(ctx, projectName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get zdb")
	}

	return zdb, nil
}

func (r *Router) deployZDB(ctx context.Context, zdb ZDB, projectName string) (ZDB, error) {
	// validate no workloads with the same name
	if err := r.validateProjectName(ctx, projectName); err != nil {
		return ZDB{}, err
	}

	// deploy
	zdbs := []workloads.ZDB{
		NewClientWorkloadFromZDB(zdb),
	}
	log.Info().Msgf("Deploying zdb: %+v", zdbs)

	clientDeployment := workloads.NewDeployment(zdb.Name, zdb.NodeID, projectName, nil, "", nil, zdbs, nil, nil)
	err := r.Client.DeploymentDeployer.Deploy(ctx, &clientDeployment)
	if err != nil {
		return ZDB{}, errors.Wrapf(err, "failed to deploy zdb with name: %s", zdb.Name)
	}

	// get the result with the computed values
	loadedZDB, err := r.Client.State.LoadZdbFromGrid(zdb.NodeID, zdb.Name, zdb.Name)
	if err != nil {
		return ZDB{}, errors.Wrapf(err, "failed to load zdb %s", zdb.Name)
	}

	result := NewZDBFromClientZDB(loadedZDB)
	result.NodeID = zdb.NodeID

	// NOTE: clean the state after deploying
	return result, nil
}

func (r *Router) deleteZDB(ctx context.Context, name string) error {
	// TODO: fix canceling
	err := r.Client.CancelByProjectName(name)
	if err != nil {
		return errors.Wrapf(err, "Failed to cancel cluster with name: %s", name)
	}

	return nil
}

func (r *Router) getZDB(ctx context.Context, projectName string) (ZDB, error) {
	// get the contract
	contracts, err := r.Client.ContractsGetter.ListContractsOfProjectName(projectName)
	if err != nil {
		return ZDB{}, errors.Wrapf(err, "failed to get contracts for project: %s", projectName)
	}

	if len(contracts.NodeContracts) != 1 {
		return ZDB{}, fmt.Errorf("contracts of project %s should be 1, but %d were found", projectName, len(contracts.NodeContracts))
	}

	contract := contracts.NodeContracts[0]

	cl, err := r.Client.NcPool.GetNodeClient(r.Client.SubstrateConn, contract.NodeID)
	if err != nil {
		return ZDB{}, errors.Wrapf(err, "failed to get client for node: %d", contract.NodeID)
	}

	cid, err := strconv.ParseUint(contract.ContractID, 10, 64)
	if err != nil {
		return ZDB{}, errors.Wrapf(err, "failed to parse contract Id: %s", contract.ContractID)
	}

	dl, err := cl.DeploymentGet(ctx, cid)
	if err != nil {
		return ZDB{}, errors.Wrapf(err, "failed to get deployment with contract Id: %s", contract.ContractID)
	}

	for _, workload := range dl.Workloads {
		if workload.Type == zos.ZDBType {
			zdb, err := workloads.NewZDBFromWorkload(&workload)
			if err != nil {
				return ZDB{}, errors.Wrapf(err, "failed to get zdb from workload: %s", workload.Name)
			}

			result := NewZDBFromClientZDB(zdb)
			result.NodeID = contract.NodeID
			return result, nil
		}
	}

	return ZDB{}, fmt.Errorf("found zdb workloads in contract %d", contract.NodeID)
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
