package procedure

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

// nodes should always be provided
func MachinesDeploy(ctx context.Context, model types.MachinesModel, client deployer.TFPluginClient) (types.MachinesModel, error) {
	// validation
	contracts, err := client.ContractsGetter.ListContractsOfProjectName(model.Name)
	if err != nil {
		return types.MachinesModel{}, errors.Wrapf(err, "failed to retreive contracts with project name %s", model.Name)
	}

	if len(contracts.NameContracts) > 0 || len(contracts.NodeContracts) > 0 || len(contracts.RentContracts) > 0 {
		return types.MachinesModel{}, fmt.Errorf("project name %s is not unique", model.Name)
	}

	// if machines don't have nodes assigned, should be assigned here

	// deploy network
	nodeList := []uint32{}
	nodeMachineMap := map[uint32][]types.Machine{}
	for _, machine := range model.Machines {
		if _, ok := nodeMachineMap[machine.NodeID]; !ok {
			nodeList = append(nodeList, machine.NodeID)
		}
		nodeMachineMap[machine.NodeID] = append(nodeMachineMap[machine.NodeID], machine)
	}

	ipRange, err := gridtypes.ParseIPNet(model.Network.IPRange)
	if err != nil {
		return types.MachinesModel{}, errors.Wrapf(err, "network ip range (%s) is invalid", model.Network.IPRange)
	}

	znet := workloads.ZNet{
		Name:        model.Network.Name,
		Nodes:       nodeList,
		IPRange:     ipRange,
		AddWGAccess: model.Network.AddWireguardAccess,
	}

	err = client.NetworkDeployer.Deploy(ctx, &znet)
	if err != nil {
		return types.MachinesModel{}, errors.Wrap(err, "failed to deploy network")
	}

	// deploy deployment
	for nodeID, machines := range nodeMachineMap {
		vms := []workloads.VM{}
		QSFSs := []workloads.QSFS{}
		disks := []workloads.Disk{}
		for _, machine := range machines {
			vm, disks, qsfss := generateVM(machine)
			qsfs := generateQSFS(machine)
			disks := generateDisk(machine)
			vms = append(vms, vm)
			//...
		}
		clientDeployment := workloads.NewDeployment(model.Name, nodeID, "", nil, model.Network.Name, disks, nil, vms, QSFSs)
		err := client.DeploymentDeployer.Deploy(ctx, &clientDeployment)

	}
	model.Network.WireguardConfig = znet.AccessWGConfig
	client.State.LoadVMFromGrid()
	/*
		- validate incoming deployment
			- project name has to be unique
			-
		- construct network deployer
			- get nodes from all machines
			- build network deployer using these nodes
		- deploy network
		- construct deployment deployer
			-
		- deploy deployment
		- construct machines model and return it
	*/

}

// machines.deelte
func MachinesDelete(ctx context.Context, name string) error

func MachineAdd(ctx context.Context, machine types.Machine, projectName string) (types.MachinesModel, error)

func MachineRemove(ctx context.Context, machineName string, projectName string) (types.MachinesModel, error)

func MachinesGet(ctx context.Context, name string) (types.MachinesModel, error)
