package router

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

// K8sCluster struct for k8s cluster
type K8sCluster struct {
	Name        string    `json:"name"`
	Master      *K8sNode  `json:"master"`
	Workers     []K8sNode `json:"workers"`
	Token       string    `json:"token"`
	NetworkName string    `json:"network_name"`
	SSHKey      string    `json:"ssh_key"`
}

// K8sNode kubernetes data
type K8sNode struct {
	Name      string `json:"name"`
	NodeID    uint32 `json:"node_id"`
	DiskSize  int    `json:"disk_size"`
	PublicIP  bool   `json:"public_ip"`
	PublicIP6 bool   `json:"public_ip6"`
	Planetary bool   `json:"planetary"`
	Flist     string `json:"flist"`
	CPU       int    `json:"cpu"`
	Memory    int    `json:"memory"`

	// computed
	ComputedIP4 string `json:"computed_ip4"`
	ComputedIP6 string `json:"computed_ip6"`
	WGIP        string `json:"wg_ip"`
	YggIP       string `json:"ygg_ip"`
}

func (r *Router) K8sDeploy(ctx context.Context, data string) (interface{}, error) {
	cluster := K8sCluster{}

	if err := json.Unmarshal([]byte(data), &cluster); err != nil {
		return K8sCluster{}, errors.Wrap(err, "failed to unmarshal k8sCluster model data")
	}

	projectName := generateProjectName(cluster.Name)

	cluster, err := r.k8sDeploy(ctx, cluster, projectName)
	if err != nil {
		return K8sCluster{}, errors.Wrap(err, "failed to deploy cluster")
	}

	return cluster, nil
}

func (r *Router) K8sDelete(ctx context.Context, data string) (interface{}, error) {
	var clusterName string

	if err := json.Unmarshal([]byte(data), &clusterName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal k8sCluster name")
	}

	projectName := generateProjectName(clusterName)

	err := r.k8sDelete(ctx, projectName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to delete cluster")
	}

	return nil, nil
}

func (r *Router) K8sGet(ctx context.Context, data string) (interface{}, error) {
	var clusterName string

	if err := json.Unmarshal([]byte(data), &clusterName); err != nil {
		return K8sCluster{}, errors.Wrap(err, "failed to unmarshal k8sCluster name")
	}

	projectName := generateProjectName(clusterName)

	cluster, err := r.k8sGet(ctx, clusterName, projectName)
	if err != nil {
		return K8sCluster{}, errors.Wrap(err, "failed to get cluster")
	}

	return cluster, nil
}

func (r *Router) k8sDeploy(ctx context.Context, cluster K8sCluster, projectName string) (K8sCluster, error) {
	// validate project name is unique
	if err := r.validateProjectName(ctx, projectName); err != nil {
		return K8sCluster{}, err
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
		return K8sCluster{}, errors.Wrapf(err, "network ip range (%s) is invalid", "10.1.0.0/16")
	}

	znet := workloads.ZNet{
		// Name:         fmt.Sprintf("%s_network", cluster.NetworkName),
		Name:         cluster.NetworkName,
		Nodes:        nodeList,
		IPRange:      ipRange,
		SolutionType: projectName,
	}

	err = r.Client.NetworkDeployer.Deploy(ctx, &znet)
	if err != nil {
		return K8sCluster{}, errors.Wrap(err, "failed to deploy network")
	}

	// map to workloads.k8sCluster
	var master workloads.K8sNode = NewClientK8sNodeFromK8sNode(*cluster.Master)
	workers := []workloads.K8sNode{}
	for _, worker := range cluster.Workers {
		workers = append(workers, NewClientK8sNodeFromK8sNode(worker))
	}

	k8s := workloads.K8sCluster{
		SolutionType: projectName,
		NetworkName:  cluster.NetworkName,
		Token:        cluster.Token,
		SSHKey:       cluster.SSHKey,
		Master:       &master,
		Workers:      workers,
	}

	// Deploy workload
	err = r.Client.K8sDeployer.Deploy(ctx, &k8s)
	if err != nil {
		return K8sCluster{}, errors.Wrapf(err, "Failed to deploy K8s Cluster")
	}

	cluster.Master.assignComputedNodeValues(*k8s.Master)
	for idx := range k8s.Workers {
		cluster.Workers[idx].assignComputedNodeValues(k8s.Workers[idx])
	}

	return cluster, nil
}

func (r *Router) k8sDelete(ctx context.Context, projectName string) error {
	err := r.Client.CancelByProjectName(projectName)
	if err != nil {
		errors.Wrapf(err, "failed to cancel project: %s", projectName)
	}

	return nil
}

func (r *Router) k8sGet(ctx context.Context, clusterName string, projectName string) (K8sCluster, error) {
	// get all contracts by project name
	contracts, err := r.Client.ContractsGetter.ListContractsOfProjectName(projectName)
	if err != nil {
		return K8sCluster{}, errors.Wrapf(err, "failed to get contracts for project: %s", projectName)
	}

	if len(contracts.NodeContracts) == 0 {
		return K8sCluster{}, fmt.Errorf("found 0 contracts for project %s", projectName)
	}

	result := K8sCluster{
		Name:    clusterName,
		Master:  &K8sNode{},
		Workers: []K8sNode{},
	}

	diskNameNodeNameMap := map[string]string{}
	nodeNameDiskSizeMap := map[string]int{}

	for _, contract := range contracts.NodeContracts {
		nodeClient, err := r.Client.NcPool.GetNodeClient(r.Client.SubstrateConn, contract.NodeID)
		if err != nil {
			return K8sCluster{}, errors.Wrapf(err, "failed to get node %d client", contract.NodeID)
		}

		contractID, err := strconv.ParseUint(contract.ContractID, 10, 64)
		if err != nil {
			return K8sCluster{}, errors.Wrapf(err, "Couldn't convert ContractID: %s", contract.ContractID)
		}

		deployment, err := nodeClient.DeploymentGet(ctx, contractID)
		if err != nil {
			return K8sCluster{}, errors.Wrapf(err, "failed to get deployment with contract id %d", contractID)
		}

		for _, workload := range deployment.Workloads {
			if workload.Type == zos.ZMachineType {
				vm := workloads.VM{}

				vm, err = workloads.NewVMFromWorkload(&workload, &deployment)
				if err != nil {
					return K8sCluster{}, errors.Wrapf(err, "Failed to get vm from workload: %s", workload)
				}

				if len(vm.Mounts) == 1 {
					diskNameNodeNameMap[vm.Mounts[0].DiskName] = vm.Name
				}

				if isWorker(vm) {
					result.Workers = append(result.Workers, NewK8sNodeFromVM(vm))
				} else {
					masterNode := NewK8sNodeFromVM(vm)
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
					return K8sCluster{}, errors.Wrapf(err, "Failed to get disk from workload: %s", workload)
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

func NewClientK8sNodeFromK8sNode(k8sNode K8sNode) workloads.K8sNode {
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

func NewK8sNodeFromVM(vm workloads.VM) K8sNode {
	return K8sNode{
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

func (k *K8sNode) assignComputedNodeValues(node workloads.K8sNode) {
	k.ComputedIP4 = node.ComputedIP
	k.ComputedIP6 = node.ComputedIP6
	k.WGIP = node.YggIP
	k.YggIP = node.IP
}

func isWorker(vm workloads.VM) bool {
	return len(vm.EnvVars["K3S_URL"]) != 0
}
