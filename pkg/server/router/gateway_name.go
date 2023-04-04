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

// GatewayNameModel struct for gateway name proxy
type GatewayNameModel struct {
	// Required
	NodeID uint32 `json:"node_id"`
	// Name the fully qualified domain name to use (cannot be present with Name)
	Name string `json:"name"`
	// Backends are list of backend ips
	Backends []zos.Backend `json:"backends"`

	// Optional
	// Passthrough whether to pass tls traffic or not
	TLSPassthrough bool   `json:"tls_passthrough"`
	Description    string `json:"description"`

	// computed

	// FQDN deployed on the node
	FQDN           string `json:"fqdn"`
	NameContractID uint64 `json:"name_contract_id"`
	ContractID     uint64 `json:"contract_id"`
}

func (r *Router) GatewayNameDeploy(ctx context.Context, data string) (interface{}, error) {
	model := GatewayNameModel{}
	if err := json.Unmarshal([]byte(data), &model); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model data")
	}

	projectName := generateProjectName(model.Name)

	res, err := r.gatewayNameDeploy(ctx, model, projectName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to deploy gateway %s", model.Name)
	}

	return res, nil
}

func (r *Router) GatewayNameGet(ctx context.Context, data string) (interface{}, error) {
	modelName := ""
	if err := json.Unmarshal([]byte(data), &modelName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model name")
	}

	projectName := generateProjectName(modelName)

	res, err := r.gatewayNameGet(ctx, projectName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to delete gateway model %s", modelName)
	}

	return res, nil
}

func (r *Router) GatewayNameDelete(ctx context.Context, data string) (interface{}, error) {
	modelName := ""
	if err := json.Unmarshal([]byte(data), &modelName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model name")
	}

	projectName := generateProjectName(modelName)

	if err := r.gatewayNameDelete(ctx, projectName); err != nil {
		return nil, errors.Wrapf(err, "failed to delete gateway model %s", modelName)
	}

	return nil, nil
}

func (r *Router) gatewayNameDeploy(ctx context.Context, gatewayNameModel GatewayNameModel, projectName string) (GatewayNameModel, error) {
	// validate that no other project is deployed with this name
	if err := r.validateProjectName(ctx, projectName); err != nil {
		return GatewayNameModel{}, err
	}

	// deploy gateway
	gateway := workloads.GatewayNameProxy{
		NodeID:         gatewayNameModel.NodeID,
		Name:           gatewayNameModel.Name,
		Backends:       gatewayNameModel.Backends,
		TLSPassthrough: gatewayNameModel.TLSPassthrough,
		Description:    gatewayNameModel.Description,
		SolutionType:   projectName,
	}

	if err := r.Client.GatewayNameDeployer.Deploy(ctx, &gateway); err != nil {
		return GatewayNameModel{}, errors.Wrapf(err, "failed to deploy gateway %s", gateway.Name)
	}

	nodeClient, err := r.Client.NcPool.GetNodeClient(r.Client.SubstrateConn, gateway.NodeID)
	if err != nil {
		return GatewayNameModel{}, errors.Wrapf(err, "failed to get node %d client", gateway.NodeID)
	}

	cfg, err := nodeClient.NetworkGetPublicConfig(ctx)
	if err != nil {
		return GatewayNameModel{}, errors.Wrapf(err, "failed to get node %d public config", gateway.NodeID)
	}

	gatewayNameModel.FQDN = fmt.Sprintf("%s.%s", gateway.Name, cfg.Domain)
	gatewayNameModel.ContractID = gateway.ContractID
	gatewayNameModel.NameContractID = gateway.NameContractID

	return gatewayNameModel, nil
}

func (r *Router) gatewayNameDelete(ctx context.Context, projectName string) error {
	if err := r.Client.CancelByProjectName(projectName); err != nil {
		return errors.Wrapf(err, "failed to cancel project %s", projectName)
	}

	return nil
}

func (r *Router) gatewayNameGet(ctx context.Context, projectName string) (GatewayNameModel, error) {
	contracts, err := r.Client.ContractsGetter.ListContractsOfProjectName(projectName)
	if err != nil {
		return GatewayNameModel{}, errors.Wrapf(err, "failed to get project %s contracts", projectName)
	}

	if len(contracts.NodeContracts) != 1 {
		return GatewayNameModel{}, fmt.Errorf("node contracts for project %s should be 1, but %d were found", projectName, len(contracts.NodeContracts))
	}

	if len(contracts.NameContracts) != 1 {
		return GatewayNameModel{}, fmt.Errorf("name contracts for project %s should be 1, but %d were found", projectName, len(contracts.NameContracts))
	}

	nodeID := contracts.NodeContracts[0].NodeID

	nodeClient, err := r.Client.NcPool.GetNodeClient(r.Client.SubstrateConn, nodeID)
	if err != nil {
		return GatewayNameModel{}, errors.Wrapf(err, "failed to get node %d client", nodeID)
	}

	nodeContractID, err := strconv.ParseUint(contracts.NodeContracts[0].ContractID, 0, 64)
	if err != nil {
		return GatewayNameModel{}, errors.Wrapf(err, "could not parse contract %s into uint64", contracts.NodeContracts[0].ContractID)
	}

	dl, err := nodeClient.DeploymentGet(ctx, nodeContractID)
	if err != nil {
		return GatewayNameModel{}, errors.Wrapf(err, "failed to get deployment with contract id %d", nodeContractID)
	}

	if len(dl.Workloads) != 1 {
		return GatewayNameModel{}, errors.Wrapf(err, "deployment should include only one gateway workload, but %d were found", len(dl.Workloads))
	}

	gatewayWorkload, err := workloads.NewGatewayNameProxyFromZosWorkload(dl.Workloads[0])
	if err != nil {
		return GatewayNameModel{}, errors.Wrapf(err, "failed to parse gateway workload data")
	}

	nameContractID, err := strconv.ParseUint(contracts.NameContracts[0].ContractID, 0, 64)
	if err != nil {
		return GatewayNameModel{}, errors.Wrapf(err, "could not parse contract %s into uint64", contracts.NameContracts[0].ContractID)
	}

	res := GatewayNameModel{
		NodeID:         nodeID,
		Name:           gatewayWorkload.Name,
		Backends:       gatewayWorkload.Backends,
		TLSPassthrough: gatewayWorkload.TLSPassthrough,
		Description:    gatewayWorkload.Description,
		FQDN:           gatewayWorkload.FQDN,
		NameContractID: nameContractID,
		ContractID:     nodeContractID,
	}

	return res, nil
}
