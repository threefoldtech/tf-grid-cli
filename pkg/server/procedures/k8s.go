package procedure

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/utils"
	"github.com/threefoldtech/zos/pkg/gridtypes"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

func K8sDeploy(ctx context.Context, cluster types.K8sCluster, client *deployer.TFPluginClient) (types.K8sCluster, error) {
	// validate project name is unique
	contracts, err := client.ContractsGetter.ListContractsOfProjectName(cluster.Name)
	if err != nil {
		return types.K8sCluster{}, errors.Wrapf(err, "Failed to retrieve contracts with project name %s", cluster.Name)
	}

	if len(contracts.NameContracts) > 0 || len(contracts.NodeContracts) > 0 || len(contracts.RentContracts) > 0 {
		return types.K8sCluster{}, fmt.Errorf("You have a cluster with the same name: %s", cluster.Name)
	}

	err = assignNodesIDsForCluster(ctx, client, &cluster)
	if err != nil {
		return types.K8sCluster{}, errors.Wrapf(err, "Couldn't find node for all cluster nodes")
	}

	// deploy network
	cluster.NetworkName = fmt.Sprintf("%s_network", cluster.Name)

	nodeList := []uint32{}
	nodeSet := map[uint32]struct{}{}
	for _, node := range cluster.Workers {
		if _, ok := nodeSet[node.NodeID]; !ok {
			nodeList = append(nodeList, node.NodeID)
			nodeSet[node.NodeID] = struct{}{}
		}
	}

	if _, ok := nodeSet[cluster.Master.NodeID]; !ok {
		nodeList = append(nodeList, cluster.Master.NodeID)
		nodeSet[cluster.Master.NodeID] = struct{}{}
	}

	ipRange, err := gridtypes.ParseIPNet("10.1.0.0/16")
	if err != nil {
		return types.K8sCluster{}, errors.Wrapf(err, "network ip range (%s) is invalid", "10.1.0.0/16")
	}

	znet := workloads.ZNet{
		// Name:         fmt.Sprintf("%s_network", cluster.NetworkName),
		Name:         cluster.NetworkName,
		Nodes:        nodeList,
		IPRange:      ipRange,
		SolutionType: cluster.Name,
	}

	err = client.NetworkDeployer.Deploy(ctx, &znet)
	if err != nil {
		return types.K8sCluster{}, errors.Wrap(err, "failed to deploy network")
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

func K8sDelete(ctx context.Context, clusterName string, client *deployer.TFPluginClient) error {
	err := client.CancelByProjectName(clusterName)
	if err != nil {
		errors.Wrapf(err, "Failed to cancel cluster with name: %s", clusterName)
	}

	return nil
}

// func K8sAddNode(ctx context.Context, clusterName string, node types.K8sNode) (types.K8sCluster, error)

// func K8sRemoveNode(ctx context.Context, clusterName string, nodeName string) (types.K8sCluster, error)

func K8sGet(ctx context.Context, clusterName string, client *deployer.TFPluginClient) (types.K8sCluster, error) {
	// get all contracts by project name
	contracts, err := client.ContractsGetter.ListContractsOfProjectName(clusterName)
	if err != nil {
		return types.K8sCluster{}, errors.Wrapf(err, "Found no clusters with this name: %s", clusterName)
	}

	result := types.K8sCluster{
		Name:    clusterName,
		Master:  &types.K8sNode{},
		Workers: []types.K8sNode{},
	}

	diskNameNodeNameMap := map[string]string{}
	nodeNameDiskSizeMap := map[string]int{}

	for _, contract := range contracts.NodeContracts {
		nodeClient, err := client.NcPool.GetNodeClient(client.SubstrateConn, contract.NodeID)

		cid, err := strconv.ParseUint(contract.ContractID, 10, 64)
		if err != nil {
			return types.K8sCluster{}, errors.Wrapf(err, "Couldn't convert ContractID: %s", contract.ContractID)
		}

		deployment, err := nodeClient.DeploymentGet(ctx, cid)

		for _, workload := range deployment.Workloads {
			if workload.Type == zos.ZMachineType {
				vm := workloads.VM{}

				vm, err = workloads.NewVMFromWorkload(&workload, &deployment)
				if err != nil {
					return types.K8sCluster{}, errors.Wrapf(err, "Failed to get vm from workload: %s", workload)
				}

				if len(vm.Mounts) == 1 {
					diskNameNodeNameMap[vm.Mounts[0].DiskName] = vm.Name
				}

				if isWorker(vm) {
					result.Workers = append(result.Workers, convertVMToK8sNode(vm))
				} else {
					masterNode := convertVMToK8sNode(vm)
					result.Master = &masterNode

					result.SSHKey = vm.EnvVars["SSH_KEY"]
					result.Token = vm.EnvVars["K3S_TOKEN"]
				}
			}
		}

		for _, workload := range deployment.Workloads {
			if workload.Type == zos.ZMountType {
				disk, err := workloads.NewDiskFromWorkload(&workload)
				if err != nil {
					return types.K8sCluster{}, errors.Wrapf(err, "Failed to get disk from workload: %s", workload)
				}
				nodeName := diskNameNodeNameMap[disk.Name]
				nodeNameDiskSizeMap[nodeName] = disk.SizeGB
			}
		}
	}

	result.Master.DiskSize = nodeNameDiskSizeMap[result.Master.Name]
	for idx := range result.Workers {
		result.Workers[idx].DiskSize = nodeNameDiskSizeMap[result.Workers[idx].Name]
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

func convertVMToK8sNode(vm workloads.VM) types.K8sNode {
	return types.K8sNode{
		Name:      vm.Name,
		PublicIP:  vm.PublicIP,
		PublicIP6: vm.PublicIP6,
		Planetary: vm.Planetary,
		Flist:     vm.Flist,
		CPU:       vm.CPU,
		Memory:    vm.Memory,

		ComputedIP4: vm.ComputedIP,
		ComputedIP6: vm.ComputedIP6,
		WGIP:        vm.IP,
		YggIP:       vm.YggIP,
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

// Assign chosen NodeIds to cluster node. with both way conversions to/from Reservations array.
func assignNodesIDsForCluster(ctx context.Context, client *deployer.TFPluginClient, cluster *types.K8sCluster) error {
	// all units unified in bytes

	workloads := []*utils.PlannedReservation{}

	if cluster.Master.NodeID == 0 {

		ms := utils.PlannedReservation{
			WorkloadName: cluster.Master.Name,
			MRU:          uint64(cluster.Master.Memory * int(gridtypes.Megabyte)),
			SRU:          uint64(cluster.Master.DiskSize * int(gridtypes.Gigabyte)),
			FarmID:       cluster.Master.FarmID,
		}

		workloads = append(workloads, &ms)
	}

	for idx := range cluster.Workers {
		if cluster.Workers[idx].NodeID == 0 {

			wr := utils.PlannedReservation{
				WorkloadName: cluster.Workers[idx].Name,
				MRU:          uint64(cluster.Workers[idx].Memory * int(gridtypes.Megabyte)),
				SRU:          uint64(cluster.Workers[idx].DiskSize * int(gridtypes.Gigabyte)),
				FarmID:       cluster.Workers[idx].FarmID,
			}

			workloads = append(workloads, &wr)
		}
	}

	err := utils.AssignNodes(ctx, client, workloads)
	if err != nil {
		return err
	}

	if cluster.Master.NodeID == 0 {
		for _, workload := range workloads {
			if workload.WorkloadName == cluster.Master.Name {
				cluster.Master.NodeID = uint32(workload.NodeID)
			}
		}
	}

	for idx := range cluster.Workers {
		if cluster.Workers[idx].NodeID == 0 {
			for _, workload := range workloads {
				if workload.WorkloadName == cluster.Workers[idx].Name {
					cluster.Workers[idx].NodeID = uint32(workload.NodeID)
				}
			}
		}
	}

	return nil
}
