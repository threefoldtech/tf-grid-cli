package procedure

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

func K8sDeploy(ctx context.Context, cluster types.K8sCluster, client deployer.TFPluginClient) (types.K8sCluster, error) {
	// validate project name is unique
	contracts, err := client.ContractsGetter.ListContractsOfProjectName(cluster.Name)
	if err != nil {
		return types.K8sCluster{}, errors.Wrapf(err, "Failed to retrieve contracts with project name %s", cluster.Name)
	}

	if len(contracts.NameContracts) > 0 || len(contracts.NodeContracts) > 0 || len(contracts.RentContracts) > 0 {
		return types.K8sCluster{}, fmt.Errorf("You have a cluster with the same name: %s", cluster.Name)
	}

	// map to workloads.k8sCluster
	var master workloads.K8sNode = convertK8sNodeToWorkload(*cluster.Master)
	workers := []workloads.K8sNode{}
	for _, worker := range cluster.Workers {
		workers = append(workers, convertK8sNodeToWorkload(worker))
	}

	k8s := workloads.K8sCluster{
		SolutionType: cluster.Name,
		NetworkName:  cluster.NetworkName,
		Token:        cluster.Token,
		SSHKey:       cluster.SSHKey,
		Master:       &master,
		Workers:      workers,
	}

	// Deploy workload
	err = client.K8sDeployer.Deploy(ctx, &k8s)
	if err != nil {
		return types.K8sCluster{}, errors.Wrapf(err, "Failed to deploy K8s Cluster")
	}

	// assign computed values to the result
	cluster.NodeDeploymentID = k8s.NodeDeploymentID
	cluster.NodesIPRange = k8s.NodesIPRange

	assignComputedNodeValues(*k8s.Master, cluster.Master)
	for idx := range k8s.Workers {
		assignComputedNodeValues(k8s.Workers[idx], &cluster.Workers[idx])
	}

	return cluster, nil
}

func K8sDelete(ctx context.Context, clusterName string, client deployer.TFPluginClient) error {
	contracts, err := client.ContractsGetter.ListContractsOfProjectName(clusterName)
	if err != nil {
		return errors.Wrapf(err, "Found no clusters with this name: %s", clusterName)
	}

	// allContracts :=  []interface{}{contracts.NameContracts, contracts.NodeContracts, contracts.RentContracts}

	for _, contract := range contracts.NodeContracts {
		num, err := strconv.ParseUint(contract.ContractID, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "Couldn't convert ContractID: %s", contract.ContractID)
		}

		err = client.SubstrateConn.CancelContract(client.Identity, num)

		if err != nil {
			return errors.Wrapf(err, "Failed deleting Contract with ContractID: %d", num)
		}
	}

	return nil
}

func K8sAddNode(ctx context.Context, clusterName string, node types.K8sNode) (types.K8sCluster, error)

func K8sRemoveNode(ctx context.Context, clusterName string, nodeName string) (types.K8sCluster, error)

func K8sGet(ctx context.Context, clusterName string, client deployer.TFPluginClient) (types.K8sCluster, error) {
	// get all contracts by project name
	contracts, err := client.ContractsGetter.ListContractsOfProjectName(clusterName)
	if err != nil {
		return types.K8sCluster{}, errors.Wrapf(err, "Found no clusters with this name: %s", clusterName)
	}

	// get deployment for each contractId. to have each {nodeId:k8sNodeName}
	masterMap := map[uint32]string{}
	workerMap := map[uint32][]string{}

	for _, contract := range contracts.NodeContracts {
		nodeClient, err := client.NcPool.GetNodeClient(client.SubstrateConn, contract.NodeID)

		cid, err := strconv.ParseUint(contracts.NodeContracts[0].ContractID, 10, 64)
		if err != nil {
			return types.K8sCluster{}, errors.Wrapf(err, "Couldn't convert ContractID: %s", contracts.NodeContracts[0].ContractID)
		}

		deployment, err := nodeClient.DeploymentGet(ctx, cid)

		for _, workload := range deployment.Workloads {
			vm := workloads.VM{}
			if workload.Type == zos.ZMachineType {
				vm, err = workloads.NewVMFromWorkload(&workload, &deployment)
				if err != nil {
					return types.K8sCluster{}, errors.Wrapf(err, "Failed to get vm from workload: %s", workload)
				}
			}
			if isWorker(vm) {
				workerMap[contract.NodeID] = append(workerMap[contract.NodeID], vm.Name)
			} else {
				masterMap[contract.NodeID] = vm.Name
			}
		}
	}

	// load the nodes from the grid
	k8s, err := client.State.LoadK8sFromGrid(
		masterMap,
		workerMap,
		clusterName,
	)

	if err != nil {
		return types.K8sCluster{}, errors.Wrapf(err, "Failed to load K8s from the grid")
	}

	// build the result
	master := convertK8sWorkloadToNode(*k8s.Master)
	workers := []types.K8sNode{}

	for _, worker := range k8s.Workers {
		workers = append(workers, convertK8sWorkloadToNode(worker))
	}

	result := types.K8sCluster{
		Name:             k8s.SolutionType,
		Token:            k8s.Token,
		SSHKey:           k8s.SSHKey,
		NodeDeploymentID: k8s.NodeDeploymentID,
		NodesIPRange:     k8s.NodesIPRange,
		Master:           &master,
	}

	return result, nil
}

func convertK8sNodeToWorkload(k8sNode types.K8sNode) workloads.K8sNode {
	return workloads.K8sNode{
		Name:      k8sNode.Name,
		Node:      k8sNode.NodeID,
		DiskSize:  k8sNode.DiskSize,
		PublicIP:  k8sNode.PublicIP,
		PublicIP6: k8sNode.PublicIP6,
		Planetary: k8sNode.Planetary,
		Flist:     k8sNode.Flist,
		CPU:       k8sNode.CPU,
		Memory:    k8sNode.Memory,
	}
}

func convertK8sWorkloadToNode(k8sNode workloads.K8sNode) types.K8sNode {
	return types.K8sNode{
		Name:      k8sNode.Name,
		NodeID:    k8sNode.Node,
		DiskSize:  k8sNode.DiskSize,
		PublicIP:  k8sNode.PublicIP,
		PublicIP6: k8sNode.PublicIP6,
		Planetary: k8sNode.Planetary,
		Flist:     k8sNode.Flist,
		CPU:       k8sNode.CPU,
		Memory:    k8sNode.Memory,

		ComputedIP4: k8sNode.ComputedIP,
		ComputedIP6: k8sNode.ComputedIP6,
		WGIP:        k8sNode.IP,
		YggIP:       k8sNode.YggIP,
	}
}

func assignComputedNodeValues(node workloads.K8sNode, resultNode *types.K8sNode) {
	resultNode.ComputedIP4 = node.ComputedIP
	resultNode.ComputedIP6 = node.ComputedIP6
	resultNode.WGIP = node.YggIP
	resultNode.YggIP = node.IP
}

func isWorker(vm workloads.VM) bool {
	return len(vm.EnvVars["K3S_URL"]) != 0
}
