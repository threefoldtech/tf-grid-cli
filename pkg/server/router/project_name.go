package router

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/graphql"
	"github.com/threefoldtech/grid3-go/workloads"
)

func (r *Router) GetContractsByProjectName(ctx context.Context, projectName string) ([]graphql.Contract, error) {
	twinContracts, err := r.Client.ContractsGetter.ListContractsByTwinID([]string{"Created, GracePeriod"})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retreive contract list with project name: %s", projectName)
	}

	contractsList := []graphql.Contract{}
	contractsList = append(contractsList, twinContracts.NameContracts...)
	contractsList = append(contractsList, twinContracts.NodeContracts...)
	contractsList = append(contractsList, twinContracts.RentContracts...)

	projectContracts := []graphql.Contract{}
	for _, contract := range contractsList {
		deploymentData, err := workloads.ParseDeploymentDate(contract.DeploymentData)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse deployment data on contarct %s", contract.ContractID)
		}

		if deploymentData.ProjectName == projectName {
			projectContracts = append(projectContracts, contract)
		}
	}

	return projectContracts, nil
}

func (r *Router) validateProjectName(ctx context.Context, projectName string) error {
	contracts, err := r.Client.ContractsGetter.ListContractsOfProjectName(projectName)
	if err != nil {
		return errors.Wrapf(err, "failed to retreive contracts with project name %s", projectName)
	}

	if len(contracts.NameContracts) > 0 || len(contracts.NodeContracts) > 0 || len(contracts.RentContracts) > 0 {
		return fmt.Errorf("invalid project name. project name (%s) is not unique", projectName)
	}

	return nil
}

func generateProjectName(projectName string) string {
	return fmt.Sprintf("%stfgridcli", projectName)
}
