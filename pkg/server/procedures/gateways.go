package procedure

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/grid3-go/graphql"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/utils"
)

func GatewayNameDeploy(ctx context.Context, gatewayNameModel types.GatewayNameModel, client *deployer.TFPluginClient) (types.GatewayNameModel, error) {
	// validate that no other project is deployed with this name
	if err := validateProjectName(ctx, gatewayNameModel.Name, client); err != nil {
		return types.GatewayNameModel{}, errors.Wrapf(err, "project name is not unique")
	}

	if gatewayNameModel.NodeID == 0 {
		nodeId, err := utils.GetGatewayNode(client)
		if err != nil {
			return types.GatewayNameModel{}, errors.Wrapf(err, "Couldn't find a gateway node")
		}

		gatewayNameModel.NodeID = nodeId
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

	if err := client.GatewayNameDeployer.Deploy(ctx, &gateway); err != nil {
		return types.GatewayNameModel{}, errors.Wrapf(err, "failed to deploy gateway %s", gateway.Name)
	}

	loadedGW, err := client.State.LoadGatewayNameFromGrid(gateway.NodeID, gateway.Name, gateway.Name)
	if err != nil {
		return types.GatewayNameModel{}, errors.Wrapf(err, "failed to load gateway %s data", gateway.Name)
	}

	gatewayNameModel.FQDN = loadedGW.FQDN
	gatewayNameModel.ContractID = gateway.ContractID
	gatewayNameModel.NameContractID = gateway.NameContractID

	return gatewayNameModel, nil
}

func GatewayNameDelete(ctx context.Context, name string, client *deployer.TFPluginClient) error {
	contractsList, err := client.ContractsGetter.ListContractsByTwinID([]string{"Created, GracePeriod"})
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

			if err := client.SubstrateConn.CancelContract(client.Identity, contractID); err != nil {
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

			if err := client.SubstrateConn.CancelContract(client.Identity, contractID); err != nil {
				return errors.Wrapf(err, "failed to cancel contract %d. retry delete", contractID)
			}

			break
		}
	}

	return nil
}

func GatewayNameGet(ctx context.Context, name string, client *deployer.TFPluginClient) (types.GatewayNameModel, error) {
	contractsList, err := client.ContractsGetter.ListContractsByTwinID([]string{"Created, GracePeriod"})
	if err != nil {
		return types.GatewayNameModel{}, errors.Wrapf(err, "failed to retreive contract list with project name: %s", name)
	}

	var nodeContract graphql.Contract
	var nameContract graphql.Contract

	for _, contract := range contractsList.NodeContracts {
		deploymentData, err := workloads.ParseDeploymentDate(contract.DeploymentData)
		if err != nil {
			return types.GatewayNameModel{}, errors.Wrapf(err, "failed to parse deployment data on contarct %s", contract.ContractID)
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

	nodeClient, err := client.NcPool.GetNodeClient(client.SubstrateConn, nodeContract.NodeID)
	if err != nil {
		return types.GatewayNameModel{}, errors.Wrapf(err, "failed to get node %d client", nodeContract.NodeID)
	}

	nodeContractID, err := strconv.ParseUint(nodeContract.ContractID, 0, 64)
	if err != nil {
		return types.GatewayNameModel{}, errors.Wrapf(err, "could not parse contract %s into uint64", nodeContract.ContractID)
	}

	dl, err := nodeClient.DeploymentGet(ctx, nodeContractID)
	if err != nil {
		return types.GatewayNameModel{}, errors.Wrapf(err, "failed to get deployment with contract id %d", nodeContractID)
	}

	gatewayWorkload := workloads.GatewayNameProxy{}

	for _, wl := range dl.Workloads {
		gatewayWorkload, err = workloads.NewGatewayNameProxyFromZosWorkload(wl)
		if err != nil {
			return types.GatewayNameModel{}, errors.Wrapf(err, "failed to parse gateway workload data")
		}

		break
	}

	nameContractID, err := strconv.ParseUint(nameContract.ContractID, 0, 64)
	if err != nil {
		return types.GatewayNameModel{}, errors.Wrapf(err, "could not parse contract %s into uint64", nameContract.ContractID)
	}

	res := types.GatewayNameModel{
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

func GatewayFQDNDeploy(ctx context.Context, gatewayFQDNModel types.GatewayFQDNModel, client *deployer.TFPluginClient) (types.GatewayFQDNModel, error) {
	if err := validateProjectName(ctx, gatewayFQDNModel.Name, client); err != nil {
		return types.GatewayFQDNModel{}, err
	}

	if gatewayFQDNModel.NodeID == 0 {
		nodeId, err := utils.GetGatewayNode(client)

		if err != nil {
			return types.GatewayFQDNModel{}, errors.Wrapf(err, "Couldn't find a gateway node")
		}

		gatewayFQDNModel.NodeID = nodeId
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

	if err := client.GatewayFQDNDeployer.Deploy(ctx, &gatewayFQDN); err != nil {
		return types.GatewayFQDNModel{}, errors.Wrapf(err, "failed to deploy gateway fqdn")
	}

	gatewayFQDNModel.ContractID = gatewayFQDN.ContractID

	return gatewayFQDNModel, nil
}

func validateProjectName(ctx context.Context, projectName string, client *deployer.TFPluginClient) error {
	twinContracts, err := client.ContractsGetter.ListContractsByTwinID([]string{"Created, GracePeriod"})
	if err != nil {
		return errors.Wrapf(err, "failed to retreive contract list with project name: %s", projectName)
	}

	contractsList := []graphql.Contract{}
	contractsList = append(contractsList, twinContracts.NameContracts...)
	contractsList = append(contractsList, twinContracts.NodeContracts...)

	for _, contract := range contractsList {
		deploymentData, err := workloads.ParseDeploymentDate(contract.DeploymentData)
		if err != nil {
			return errors.Wrapf(err, "failed to parse deployment data on contarct %s", contract.ContractID)
		}

		if deploymentData.ProjectName == projectName {
			return errors.Wrapf(err, "model name %s is not unique. cancel any contracts with this project name first", projectName)
		}
	}

	return nil
}

func GatewayFQDNDelete(ctx context.Context, name string, client *deployer.TFPluginClient) error {
	if err := client.CancelByProjectName(name); err != nil {
		return errors.Wrapf(err, "failed to delete gateway fqdn model contracts")
	}

	return nil
}

func GatewayFQDNGet(ctx context.Context, name string, client *deployer.TFPluginClient) (types.GatewayFQDNModel, error) {
	contractsList, err := client.ContractsGetter.ListContractsByTwinID([]string{"Created, GracePeriod"})
	if err != nil {
		return types.GatewayFQDNModel{}, errors.Wrapf(err, "failed to retreive contract list with project name: %s", name)
	}

	var nodeContract graphql.Contract

	for _, contract := range contractsList.NodeContracts {
		deploymentData, err := workloads.ParseDeploymentDate(contract.DeploymentData)
		if err != nil {
			return types.GatewayFQDNModel{}, errors.Wrapf(err, "failed to parse deployment data on contarct %s", contract.ContractID)
		}

		if deploymentData.ProjectName == name {
			nodeContract = contract
			break
		}
	}

	nodeClient, err := client.NcPool.GetNodeClient(client.SubstrateConn, nodeContract.NodeID)
	if err != nil {
		return types.GatewayFQDNModel{}, errors.Wrapf(err, "failed to get node %d client", nodeContract.NodeID)
	}

	nodeContractID, err := strconv.ParseUint(nodeContract.ContractID, 0, 64)
	if err != nil {
		return types.GatewayFQDNModel{}, errors.Wrapf(err, "could not parse contract %s into uint64", nodeContract.ContractID)
	}

	dl, err := nodeClient.DeploymentGet(ctx, nodeContractID)
	if err != nil {
		return types.GatewayFQDNModel{}, errors.Wrapf(err, "failed to get deployment with contract id %d", nodeContractID)
	}

	gatewayWorkload := workloads.GatewayFQDNProxy{}

	for _, wl := range dl.Workloads {
		gatewayWorkload, err = workloads.NewGatewayFQDNProxyFromZosWorkload(wl)
		if err != nil {
			return types.GatewayFQDNModel{}, errors.Wrapf(err, "failed to parse gateway workload data")
		}

		break
	}

	res := types.GatewayFQDNModel{
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
