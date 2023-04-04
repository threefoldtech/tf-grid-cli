package router

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

// GatewayFQDNModel for gateway FQDN proxy
type GatewayFQDNModel struct {
	// required
	NodeID uint32 `json:"node_id"`
	// Backends are list of backend ips
	Backends []zos.Backend `json:"backends"`
	// FQDN deployed on the node
	FQDN string `json:"fqdn"`
	// Name is the workload name
	Name string `json:"name"`

	// optional
	// Passthrough whether to pass tls traffic or not
	TLSPassthrough bool   `json:"tls_passthrough"`
	Description    string `json:"description"`

	// computed
	ContractID uint64 `json:"contract_id"`
}

func (r *Router) GatewayFQDNDeploy(ctx context.Context, data string) (interface{}, error) {
	model := GatewayFQDNModel{}
	if err := json.Unmarshal([]byte(data), &model); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model data")
	}

	projectName := generateProjectName(model.Name)

	res, err := r.gatewayFQDNDeploy(ctx, model, projectName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to deploy gateway %s", model.Name)
	}

	return res, nil
}

func (r *Router) GatewayFQDNGet(ctx context.Context, data string) (interface{}, error) {
	modelName := ""
	if err := json.Unmarshal([]byte(data), &modelName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model name")
	}

	projectName := generateProjectName(modelName)

	res, err := r.gatewayFQDNGet(ctx, projectName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to delete gateway model %s", modelName)
	}

	return res, nil
}

func (r *Router) GatewayFQDNDelete(ctx context.Context, data string) (interface{}, error) {
	modelName := ""
	if err := json.Unmarshal([]byte(data), &modelName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model name")
	}

	projectName := generateProjectName(modelName)

	if err := r.gatewayFQDNDelete(ctx, projectName); err != nil {
		return nil, errors.Wrapf(err, "failed to delete gateway model %s", modelName)
	}

	return nil, nil
}

func (r *Router) gatewayFQDNDeploy(ctx context.Context, gatewayFQDNModel GatewayFQDNModel, projectName string) (GatewayFQDNModel, error) {
	if err := r.validateProjectName(ctx, gatewayFQDNModel.Name); err != nil {
		return GatewayFQDNModel{}, err
	}

	gatewayFQDN := workloads.GatewayFQDNProxy{
		NodeID:         gatewayFQDNModel.NodeID,
		Backends:       gatewayFQDNModel.Backends,
		FQDN:           gatewayFQDNModel.FQDN,
		Name:           gatewayFQDNModel.Name,
		TLSPassthrough: gatewayFQDNModel.TLSPassthrough,
		Description:    gatewayFQDNModel.Description,
		SolutionType:   projectName,
	}

	if err := r.Client.GatewayFQDNDeployer.Deploy(ctx, &gatewayFQDN); err != nil {
		return GatewayFQDNModel{}, errors.Wrapf(err, "failed to deploy gateway fqdn")
	}

	gatewayFQDNModel.ContractID = gatewayFQDN.ContractID

	return gatewayFQDNModel, nil
}

func (r *Router) gatewayFQDNDelete(ctx context.Context, projectName string) error {
	if err := r.Client.CancelByProjectName(projectName); err != nil {
		return errors.Wrapf(err, "failed to delete gateway fqdn model contracts")
	}

	return nil
}

func (r *Router) gatewayFQDNGet(ctx context.Context, projectName string) (GatewayFQDNModel, error) {
	contracts, err := r.Client.ContractsGetter.ListContractsOfProjectName(projectName)
	if err != nil {
		return GatewayFQDNModel{}, errors.Wrapf(err, "failed to get project %s contracts", projectName)
	}

	if len(contracts.NodeContracts) != 1 {
		return GatewayFQDNModel{}, fmt.Errorf("node contracts for project %s should be 1, but %d were found", projectName, len(contracts.NodeContracts))
	}

	nodeID := contracts.NodeContracts[0].NodeID

	nodeClient, err := r.Client.NcPool.GetNodeClient(r.Client.SubstrateConn, nodeID)
	if err != nil {
		return GatewayFQDNModel{}, errors.Wrapf(err, "failed to get node %d client", nodeID)
	}

	nodeContractID, err := strconv.ParseUint(contracts.NodeContracts[0].ContractID, 0, 64)
	if err != nil {
		return GatewayFQDNModel{}, errors.Wrapf(err, "could not parse contract %s into uint64", contracts.NodeContracts[0].ContractID)
	}

	dl, err := nodeClient.DeploymentGet(ctx, nodeContractID)
	if err != nil {
		return GatewayFQDNModel{}, errors.Wrapf(err, "failed to get deployment with contract id %d", nodeContractID)
	}

	gatewayWorkload := workloads.GatewayFQDNProxy{}

	for _, wl := range dl.Workloads {
		gatewayWorkload, err = workloads.NewGatewayFQDNProxyFromZosWorkload(wl)
		if err != nil {
			return GatewayFQDNModel{}, errors.Wrapf(err, "failed to parse gateway workload data")
		}

		break
	}

	res := GatewayFQDNModel{
		NodeID:         nodeID,
		Name:           gatewayWorkload.Name,
		Backends:       gatewayWorkload.Backends,
		TLSPassthrough: gatewayWorkload.TLSPassthrough,
		Description:    gatewayWorkload.Description,
		FQDN:           gatewayWorkload.FQDN,
		ContractID:     nodeContractID,
	}

	return res, nil
}
