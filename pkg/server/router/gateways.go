package router

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/graphql"
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

	TLSPassthrough bool   `json:"tls_passthrough"`
	Description    string `json:"description"`
	// Optional
	// Passthrough whether to pass tls traffic or not

	// computed

	// FQDN deployed on the node
	// NodeDeploymentID map[uint32]uint64
	FQDN           string `json:"fqdn"`
	NameContractID uint64 `json:"name_contract_id"`
	ContractID     uint64 `json:"contract_id"`
}

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

	// SolutionType     string
	// NodeDeploymentID map[uint32]uint64

	// computed
	ContractID uint64 `json:"contract_id"`
}

func (r *Router) GatewayNameDeploy(ctx context.Context, data string) (interface{}, error) {
	model := GatewayNameModel{}
	if err := json.Unmarshal([]byte(data), &model); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model data")
	}

	originalProjectName := model.Name
	cliProjectName := generateProjectName(model.Name)
	model.Name = cliProjectName

	res, err := r.gatewayNameDeploy(ctx, model)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to deploy gateway %s", model.Name)
	}

	res.Name = originalProjectName

	return res, nil
}

func (r *Router) GatewayNameGet(ctx context.Context, data string) (interface{}, error) {
	modelName := ""
	if err := json.Unmarshal([]byte(data), &modelName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model name")
	}

	originalProjectName := modelName
	cliProjectName := generateProjectName(modelName)
	modelName = cliProjectName

	res, err := r.gatewayNameGet(ctx, modelName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to delete gateway model %s", modelName)
	}

	res.Name = originalProjectName

	return res, nil
}

func (r *Router) GatewayNameDelete(ctx context.Context, data string) (interface{}, error) {
	modelName := ""
	if err := json.Unmarshal([]byte(data), &modelName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model name")
	}

	cliProjectName := generateProjectName(modelName)
	modelName = cliProjectName

	if err := r.gatewayNameDelete(ctx, modelName); err != nil {
		return nil, errors.Wrapf(err, "failed to delete gateway model %s", modelName)
	}

	return nil, nil
}

func (r *Router) GatewayFQDNDeploy(ctx context.Context, data string) (interface{}, error) {
	model := GatewayFQDNModel{}
	if err := json.Unmarshal([]byte(data), &model); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model data")
	}

	originalProjectName := model.Name
	cliProjectName := generateProjectName(model.Name)
	model.Name = cliProjectName

	res, err := r.gatewayFQDNDeploy(ctx, model)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to deploy gateway %s", model.Name)
	}

	model.Name = originalProjectName

	return res, nil
}

func (r *Router) GatewayFQDNGet(ctx context.Context, data string) (interface{}, error) {
	modelName := ""
	if err := json.Unmarshal([]byte(data), &modelName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model name")
	}

	originalProjectName := modelName
	cliProjectName := generateProjectName(modelName)
	modelName = cliProjectName

	res, err := r.gatewayFQDNGet(ctx, modelName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to delete gateway model %s", modelName)
	}

	res.Name = originalProjectName

	return res, nil
}

func (r *Router) GatewayFQDNDelete(ctx context.Context, data string) (interface{}, error) {
	modelName := ""
	if err := json.Unmarshal([]byte(data), &modelName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal gateway model name")
	}

	cliProjectName := generateProjectName(modelName)
	modelName = cliProjectName

	if err := r.gatewayFQDNDelete(ctx, modelName); err != nil {
		return nil, errors.Wrapf(err, "failed to delete gateway model %s", modelName)
	}

	return nil, nil
}

func (r *Router) gatewayNameDeploy(ctx context.Context, gatewayNameModel GatewayNameModel) (GatewayNameModel, error) {
	// validate that no other project is deployed with this name
	if err := r.validateProjectName(ctx, gatewayNameModel.Name); err != nil {
		return GatewayNameModel{}, err
	}

	// deploy gateway
	gateway := workloads.GatewayNameProxy{
		NodeID:         gatewayNameModel.NodeID,
		Name:           gatewayNameModel.Name,
		Backends:       gatewayNameModel.Backends,
		TLSPassthrough: gatewayNameModel.TLSPassthrough,
		Description:    gatewayNameModel.Description,
		SolutionType:   gatewayNameModel.Name,
	}

	if err := r.Client.GatewayNameDeployer.Deploy(ctx, &gateway); err != nil {
		return GatewayNameModel{}, errors.Wrapf(err, "failed to deploy gateway %s", gateway.Name)
	}

	loadedGW, err := r.Client.State.LoadGatewayNameFromGrid(gateway.NodeID, gateway.Name, gateway.Name)
	if err != nil {
		return GatewayNameModel{}, errors.Wrapf(err, "failed to load gateway %s data", gateway.Name)
	}

	gatewayNameModel.FQDN = loadedGW.FQDN
	gatewayNameModel.ContractID = gateway.ContractID
	gatewayNameModel.NameContractID = gateway.NameContractID

	return gatewayNameModel, nil
}

func (r *Router) gatewayNameDelete(ctx context.Context, name string) error {
	contractsList, err := r.Client.ContractsGetter.ListContractsByTwinID([]string{"Created, GracePeriod"})
	if err != nil {
		return errors.Wrapf(err, "failed to retreive contract list with project name: %s", name)
	}

	for _, contract := range contractsList.NodeContracts {
		deploymentData, err := workloads.ParseDeploymentDate(contract.DeploymentData)
		if err != nil {
			return errors.Wrapf(err, "failed to parse deployment data on contarct %s", contract.ContractID)
		}

		if deploymentData.ProjectName == name {
			contractID, err := strconv.ParseUint(contract.ContractID, 0, 64)
			if err != nil {
				return errors.Wrapf(err, "could not parse contract %s into uint64", contract.ContractID)
			}

			if err := r.Client.SubstrateConn.CancelContract(r.Client.Identity, contractID); err != nil {
				return errors.Wrapf(err, "failed to cancel contract %d. retry delete", contractID)
			}
		}
	}

	for _, contract := range contractsList.NameContracts {
		if contract.Name == name {
			contractID, err := strconv.ParseUint(contract.ContractID, 0, 64)
			if err != nil {
				return errors.Wrapf(err, "could not parse contract %s into uint64", contract.ContractID)
			}

			if err := r.Client.SubstrateConn.CancelContract(r.Client.Identity, contractID); err != nil {
				return errors.Wrapf(err, "failed to cancel contract %d. retry delete", contractID)
			}

			break
		}
	}

	return nil
}

func (r *Router) gatewayNameGet(ctx context.Context, name string) (GatewayNameModel, error) {
	contractsList, err := r.Client.ContractsGetter.ListContractsByTwinID([]string{"Created, GracePeriod"})
	if err != nil {
		return GatewayNameModel{}, errors.Wrapf(err, "failed to retreive contract list with project name: %s", name)
	}

	var nodeContract graphql.Contract
	var nameContract graphql.Contract

	for _, contract := range contractsList.NodeContracts {
		deploymentData, err := workloads.ParseDeploymentDate(contract.DeploymentData)
		if err != nil {
			return GatewayNameModel{}, errors.Wrapf(err, "failed to parse deployment data on contarct %s", contract.ContractID)
		}

		if deploymentData.ProjectName == name {
			nodeContract = contract
			break
		}
	}

	for _, contract := range contractsList.NameContracts {
		if contract.Name == name {
			nameContract = contract
			break
		}
	}

	nodeClient, err := r.Client.NcPool.GetNodeClient(r.Client.SubstrateConn, nodeContract.NodeID)
	if err != nil {
		return GatewayNameModel{}, errors.Wrapf(err, "failed to get node %d client", nodeContract.NodeID)
	}

	nodeContractID, err := strconv.ParseUint(nodeContract.ContractID, 0, 64)
	if err != nil {
		return GatewayNameModel{}, errors.Wrapf(err, "could not parse contract %s into uint64", nodeContract.ContractID)
	}

	dl, err := nodeClient.DeploymentGet(ctx, nodeContractID)
	if err != nil {
		return GatewayNameModel{}, errors.Wrapf(err, "failed to get deployment with contract id %d", nodeContractID)
	}

	gatewayWorkload := workloads.GatewayNameProxy{}

	for _, wl := range dl.Workloads {
		gatewayWorkload, err = workloads.NewGatewayNameProxyFromZosWorkload(wl)
		if err != nil {
			return GatewayNameModel{}, errors.Wrapf(err, "failed to parse gateway workload data")
		}

		break
	}

	nameContractID, err := strconv.ParseUint(nameContract.ContractID, 0, 64)
	if err != nil {
		return GatewayNameModel{}, errors.Wrapf(err, "could not parse contract %s into uint64", nameContract.ContractID)
	}

	res := GatewayNameModel{
		NodeID:         nodeContract.NodeID,
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

func (r *Router) gatewayFQDNDeploy(ctx context.Context, gatewayFQDNModel GatewayFQDNModel) (GatewayFQDNModel, error) {
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
		SolutionType:   gatewayFQDNModel.Name,
	}

	if err := r.Client.GatewayFQDNDeployer.Deploy(ctx, &gatewayFQDN); err != nil {
		return GatewayFQDNModel{}, errors.Wrapf(err, "failed to deploy gateway fqdn")
	}

	gatewayFQDNModel.ContractID = gatewayFQDN.ContractID

	return gatewayFQDNModel, nil
}

func (r *Router) gatewayFQDNDelete(ctx context.Context, name string) error {
	if err := r.Client.CancelByProjectName(name); err != nil {
		return errors.Wrapf(err, "failed to delete gateway fqdn model contracts")
	}

	return nil
}

func (r *Router) gatewayFQDNGet(ctx context.Context, name string) (GatewayFQDNModel, error) {
	contractsList, err := r.Client.ContractsGetter.ListContractsByTwinID([]string{"Created, GracePeriod"})
	if err != nil {
		return GatewayFQDNModel{}, errors.Wrapf(err, "failed to retreive contract list with project name: %s", name)
	}

	var nodeContract graphql.Contract

	for _, contract := range contractsList.NodeContracts {
		deploymentData, err := workloads.ParseDeploymentDate(contract.DeploymentData)
		if err != nil {
			return GatewayFQDNModel{}, errors.Wrapf(err, "failed to parse deployment data on contarct %s", contract.ContractID)
		}

		if deploymentData.ProjectName == name {
			nodeContract = contract
			break
		}
	}

	nodeClient, err := r.Client.NcPool.GetNodeClient(r.Client.SubstrateConn, nodeContract.NodeID)
	if err != nil {
		return GatewayFQDNModel{}, errors.Wrapf(err, "failed to get node %d client", nodeContract.NodeID)
	}

	nodeContractID, err := strconv.ParseUint(nodeContract.ContractID, 0, 64)
	if err != nil {
		return GatewayFQDNModel{}, errors.Wrapf(err, "could not parse contract %s into uint64", nodeContract.ContractID)
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
		NodeID:         nodeContract.NodeID,
		Name:           gatewayWorkload.Name,
		Backends:       gatewayWorkload.Backends,
		TLSPassthrough: gatewayWorkload.TLSPassthrough,
		Description:    gatewayWorkload.Description,
		FQDN:           gatewayWorkload.FQDN,
		ContractID:     nodeContractID,
	}

	return res, nil
}
