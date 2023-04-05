package utils

import (
	"fmt"
	"math/rand"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/deployer"
	proxyTypes "github.com/threefoldtech/grid_proxy_server/pkg/types"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

var (
	status  = "up"
	trueVal = true
)

type Reservation struct {
	SelectedNodeID uint64
	WorkloadName   string
	NeededMemory   uint64
	NeededStorage  uint64
}

// Assign chosen NodeIds to cluster node. with both way conversions to/from Reservations array.
func AssignNodesIDsForCluster(client *deployer.TFPluginClient, cluster *types.K8sCluster) error {
	// all units unified in bytes

	workloads := []*Reservation{}

	if cluster.Master.NodeID == 0 {

		ms := Reservation{
			WorkloadName:  cluster.Master.Name,
			NeededMemory:  uint64(cluster.Master.Memory * int(gridtypes.Megabyte)),
			NeededStorage: uint64(cluster.Master.DiskSize * int(gridtypes.Gigabyte)),
		}
		workloads = append(workloads, &ms)
	}

	for idx := range cluster.Workers {
		if cluster.Workers[idx].NodeID == 0 {
			workloads = append(workloads, &Reservation{
				WorkloadName:  cluster.Workers[idx].Name,
				NeededMemory:  uint64(cluster.Workers[idx].Memory * int(gridtypes.Megabyte)),
				NeededStorage: uint64(cluster.Workers[idx].DiskSize * int(gridtypes.Gigabyte)),
			})
		}
	}

	err := getNodes(client, workloads)
	if err != nil {
		return err
	}

	if cluster.Master.NodeID == 0 {
		for _, workload := range workloads {
			if workload.WorkloadName == cluster.Master.Name {
				if workload.SelectedNodeID == 0 {
					return fmt.Errorf("Couldn't find node for workload")
				}
				cluster.Master.NodeID = uint32(workload.SelectedNodeID)
			}
		}
	}

	for idx := range cluster.Workers {
		if cluster.Workers[idx].NodeID == 0 {
			for _, workload := range workloads {
				if workload.WorkloadName == cluster.Workers[idx].Name {
					if workload.SelectedNodeID == 0 {
						return fmt.Errorf("Couldn't find node for workload")
					}
					cluster.Workers[idx].NodeID = uint32(workload.SelectedNodeID)
				}
			}
		}
	}

	return nil
}

// Assign chosen NodeIds to machines vm. with both way conversions to/from Reservations array.
func AssignNodesIDsForMachines(client *deployer.TFPluginClient, machines *types.MachinesModel) error {
	// all units unified in bytes

	workloads := []*Reservation{}

	for idx := range machines.Machines {
		if machines.Machines[idx].NodeID == 0 {
			neededStorage := 0
			for _, disk := range machines.Machines[idx].Disks {
				neededStorage += disk.SizeGB * int(gridtypes.Gigabyte)
			}
			for _, qsfs := range machines.Machines[idx].QSFSs {
				neededStorage += int(qsfs.MaxZDBDataDirSize) * int(gridtypes.Gigabyte)
			}
			neededStorage += machines.Machines[idx].RootfsSize * int(gridtypes.Megabyte)

			workloads = append(workloads, &Reservation{
				WorkloadName:  machines.Machines[idx].Name,
				NeededMemory:  uint64(machines.Machines[idx].Memory * int(gridtypes.Megabyte)),
				NeededStorage: uint64(neededStorage),
			})
		}
	}

	err := getNodes(client, workloads)
	if err != nil {
		return err
	}

	for idx := range machines.Machines {
		if machines.Machines[idx].NodeID == 0 {
			for _, workload := range workloads {
				if workload.WorkloadName == machines.Machines[idx].Name {
					if workload.SelectedNodeID == 0 {
						return fmt.Errorf("Couldn't find node for workload")
					}
					machines.Machines[idx].NodeID = uint32(workload.SelectedNodeID)
				}
			}
		}
	}

	return nil
}

func getFreeResources(node proxyTypes.Node) proxyTypes.Capacity {
	return proxyTypes.Capacity{
		MRU: (node.TotalResources.MRU - node.UsedResources.MRU),
		SRU: (node.TotalResources.SRU - node.UsedResources.SRU),
	}
}

func checkNodeAvailability(node proxyTypes.Node, workload Reservation, reservedCapacity map[uint64]Reservation) bool {
	if uint64(getFreeResources(node).MRU) >= reservedCapacity[uint64(node.NodeID)].NeededMemory+uint64(workload.NeededMemory) &&
		uint64(getFreeResources(node).SRU) >= reservedCapacity[uint64(node.NodeID)].NeededStorage+uint64(workload.NeededStorage) {
		return true
	}
	return false
}

// Searching for node for each workload considering the reserved capacity by workloads in the same deployment.
// Assign the NodeID if found one or return it with NodeID: 0
func getNodes(client *deployer.TFPluginClient, workloads []*Reservation) error {
	reservedCapacity := make(map[uint64]Reservation)

	for _, workload := range workloads {
		mru := uint64(workload.NeededMemory)
		sru := uint64(workload.NeededStorage)
		ips := uint64(0)

		options := proxyTypes.NodeFilter{
			Status:  &status,
			FreeMRU: &mru,
			FreeSRU: &sru,
			FreeIPs: &ips,
		}

		nodes, err := deployer.FilterNodes(client.GridProxyClient, options)
		if err != nil || len(nodes) == 0 {
			errors.Wrapf(err, "Couldn't find node for the provided filters: %+v", options)
		}

		selectedNode := proxyTypes.Node{}

		for _, node := range nodes {
			valid := checkNodeAvailability(node, *workload, reservedCapacity)
			if valid {
				selectedNode = node
				break
			}
		}

		reservedCapacity[uint64(selectedNode.NodeID)] = Reservation{
			NeededMemory:  reservedCapacity[uint64(selectedNode.NodeID)].NeededMemory + workload.NeededMemory,
			NeededStorage: reservedCapacity[uint64(selectedNode.NodeID)].NeededStorage + workload.NeededStorage,
		}

		workload.SelectedNodeID = uint64(selectedNode.NodeID)
	}

	return nil
}

func GetGatewayNode(client *deployer.TFPluginClient) (uint32, error) {
	options := proxyTypes.NodeFilter{
		Status: &status,
		IPv4:   &trueVal,
		Domain: &trueVal,
	}

	nodes, err := deployer.FilterNodes(client.GridProxyClient, options)
	if err != nil || len(nodes) == 0 {
		return 0, errors.Wrapf(err, "Couldn't find node for the provided filters: %+v", options)
	}

	return uint32(nodes[rand.Intn(len(nodes))].NodeID), nil
}

func GetNodeForZdb(client *deployer.TFPluginClient, size uint64) (uint32, error) {
	options := proxyTypes.NodeFilter{
		Status:  &status,
		FreeSRU: &size,
	}

	nodes, err := deployer.FilterNodes(client.GridProxyClient, options)
	if err != nil || len(nodes) == 0 {
		return 0, errors.Wrapf(err, "Couldn't find node for the provided filters: %+v", options)
	}

	return uint32(nodes[rand.Intn(len(nodes))].NodeID), nil
}
