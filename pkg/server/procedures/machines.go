package procedure

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
	"github.com/threefoldtech/zos/pkg/gridtypes"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// nodes should always be provided
func MachinesDeploy(ctx context.Context, model types.MachinesModel, client *deployer.TFPluginClient) (types.MachinesModel, error) {
	/*
		- validate incoming deployment
			- project name has to be unique
		- construct network deployer
			- get nodes from all machines
			- build network deployer using these nodes
		- deploy network
		- construct deployment deployer
		- deploy deployment
		- construct machines model and return it
	*/

	// validation
	contracts, err := client.ContractsGetter.ListContractsOfProjectName(model.Name)
	if err != nil {
		return types.MachinesModel{}, errors.Wrapf(err, "failed to retreive contracts with project name %s", model.Name)
	}

	if len(contracts.NameContracts) > 0 || len(contracts.NodeContracts) > 0 || len(contracts.RentContracts) > 0 {
		return types.MachinesModel{}, fmt.Errorf("project name %s is not unique", model.Name)
	}

	// TODO: if machines don't have nodes assigned, should be assigned here

	// deploy network
	znet, err := deployNetwork(ctx, &model, client)
	if err != nil {
		return types.MachinesModel{}, err
	}

	// deploy deployment
	if err := deployDeployment(ctx, &model, client); err != nil {
		// TODO: if error happens midway, all created contracts should be deleted
		return types.MachinesModel{}, err
	}

	// construct result
	if err := constructResult(&model, znet, client); err != nil {
		return types.MachinesModel{}, err
	}

	return model, nil
}

func constructResult(model *types.MachinesModel, znet *workloads.ZNet, client *deployer.TFPluginClient) error {
	model.Network.WireguardConfig = znet.AccessWGConfig

	for idx, m := range model.Machines {
		vm, err := client.State.LoadVMFromGrid(m.NodeID, m.Name, model.Name)
		if err != nil {
			return errors.Wrap(err, "deployment was successful, but failed to construct result")
		}

		// get machine ips
		model.Machines[idx].ComputedIP4 = vm.ComputedIP
		model.Machines[idx].ComputedIP6 = vm.ComputedIP6
		model.Machines[idx].YggIP = vm.YggIP
		model.Machines[idx].WGIP = vm.IP

		for idy, qsfs := range model.Machines[idx].QSFSs {
			q, err := client.State.LoadQSFSFromGrid(m.NodeID, qsfs.Name, model.Name)
			if err != nil {
				return errors.Wrap(err, "deployment was successful, but failed to construct result")
			}
			model.Machines[idx].QSFSs[idy].MetricsEndpoint = q.MetricsEndpoint
		}
	}

	return nil
}

func deployDeployment(ctx context.Context, model *types.MachinesModel, client *deployer.TFPluginClient) error {
	nodeMachineMap := map[uint32][]*types.Machine{}
	for idx, machine := range model.Machines {
		nodeMachineMap[machine.NodeID] = append(nodeMachineMap[machine.NodeID], &model.Machines[idx])
	}

	networkName := fmt.Sprintf("%s_network", model.Name)

	for nodeID, machines := range nodeMachineMap {
		vms := []workloads.VM{}
		QSFSs := []workloads.QSFS{}
		disks := []workloads.Disk{}

		for _, machine := range machines {
			nodeVM, nodeDisks, nodeQSFSs := extractWorkloads(machine, networkName)
			vms = append(vms, nodeVM)
			QSFSs = append(QSFSs, nodeQSFSs...)
			disks = append(disks, nodeDisks...)
		}

		clientDeployment := workloads.NewDeployment(model.Name, nodeID, model.Name, nil, networkName, disks, nil, vms, QSFSs)
		if err := client.DeploymentDeployer.Deploy(ctx, &clientDeployment); err != nil {
			return errors.Wrap(err, "failed to deploy")
		}
	}

	return nil
}

func deployNetwork(ctx context.Context, model *types.MachinesModel, client *deployer.TFPluginClient) (*workloads.ZNet, error) {
	nodeList := []uint32{}
	nodeSet := map[uint32]struct{}{}
	for _, machine := range model.Machines {
		if _, ok := nodeSet[machine.NodeID]; !ok {
			nodeList = append(nodeList, machine.NodeID)
			nodeSet[machine.NodeID] = struct{}{}
		}
	}

	ipRange, err := gridtypes.ParseIPNet(model.Network.IPRange)
	if err != nil {
		return nil, errors.Wrapf(err, "network ip range (%s) is invalid", model.Network.IPRange)
	}

	znet := workloads.ZNet{
		Name:         fmt.Sprintf("%s_network", model.Name),
		Nodes:        nodeList,
		IPRange:      ipRange,
		AddWGAccess:  model.Network.AddWireguardAccess,
		SolutionType: model.Name,
	}

	if znet.AddWGAccess == true {
		privateKey, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate wireguard private key")
		}
		znet.ExternalSK = privateKey
	}

	err = client.NetworkDeployer.Deploy(ctx, &znet)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deploy network")
	}

	return &znet, nil
}

func extractWorkloads(machine *types.Machine, networkName string) (workloads.VM, []workloads.Disk, []workloads.QSFS) {
	disks := []workloads.Disk{}
	qsfss := []workloads.QSFS{}
	mounts := []workloads.Mount{}
	zlogs := []workloads.Zlog{}

	for idx, disk := range machine.Disks {
		diskName := fmt.Sprintf("%s_disk_%d", machine.Name, idx)
		disks = append(disks, workloads.Disk{
			Name:        diskName,
			SizeGB:      disk.SizeGB,
			Description: disk.Description,
		})
		mounts = append(mounts, workloads.Mount{
			DiskName:   diskName,
			MountPoint: disk.MountPoint,
		})
	}

	for idx, qsfs := range machine.QSFSs {
		metaBackends := []workloads.Backend{}
		for _, b := range qsfs.Metadata.Backends {
			metaBackends = append(metaBackends, workloads.Backend{
				Address:   b.Address,
				Namespace: b.Namespace,
				Password:  b.Password,
			})
		}
		groups := []workloads.Group{}
		for _, group := range qsfs.Groups {
			bs := workloads.Backends{}
			for _, b := range group.Backends {
				bs = append(bs, workloads.Backend{
					Address:   b.Address,
					Namespace: b.Namespace,
					Password:  b.Password,
				})
			}
			groups = append(groups, workloads.Group{Backends: bs})
		}

		qsfss = append(qsfss, workloads.QSFS{
			Name:                 fmt.Sprintf("%s_qsfs_%d", machine.Name, idx),
			Description:          qsfs.Description,
			Cache:                qsfs.Cache,
			MinimalShards:        qsfs.MinimalShards,
			ExpectedShards:       qsfs.ExpectedShards,
			RedundantGroups:      qsfs.RedundantGroups,
			RedundantNodes:       qsfs.RedundantNodes,
			MaxZDBDataDirSize:    qsfs.MaxZDBDataDirSize,
			EncryptionAlgorithm:  qsfs.EncryptionAlgorithm,
			EncryptionKey:        qsfs.EncryptionKey,
			CompressionAlgorithm: qsfs.CompressionAlgorithm,
			Metadata: workloads.Metadata{
				Type:                qsfs.Metadata.Type,
				Prefix:              qsfs.Metadata.Prefix,
				EncryptionAlgorithm: qsfs.Metadata.EncryptionAlgorithm,
				EncryptionKey:       qsfs.Metadata.EncryptionKey,
				Backends:            metaBackends,
			},
			Groups:          groups,
			MetricsEndpoint: qsfs.MetricsEndpoint,
		})
	}

	for _, zlog := range machine.Zlogs {
		zlogs = append(zlogs, workloads.Zlog{
			Zmachine: machine.Name,
			Output:   zlog.Output,
		})
	}

	vm := workloads.VM{
		Name:        machine.Name,
		Flist:       machine.Flist,
		PublicIP:    machine.PublicIP,
		PublicIP6:   machine.PublicIP6,
		Planetary:   machine.Planetary,
		Description: machine.Description,
		CPU:         machine.CPU,
		Memory:      machine.Memory,
		RootfsSize:  machine.RootfsSize,
		Entrypoint:  machine.Entrypoint,
		Mounts:      mounts,
		Zlogs:       zlogs,
		EnvVars:     machine.EnvVars,
		NetworkName: networkName,
	}

	return vm, disks, qsfss
}

func MachinesDelete(ctx context.Context, name string, client *deployer.TFPluginClient) error {
	if err := client.CancelByProjectName(name); err != nil {
		return errors.Wrapf(err, "failed to cancel contracts")
	}

	return nil
}

func MachinesGet(ctx context.Context, name string, client *deployer.TFPluginClient) (types.MachinesModel, error) {
	model := types.MachinesModel{}

	contracts, err := client.ContractsGetter.ListContractsOfProjectName(name)
	if err != nil {
		return types.MachinesModel{}, errors.Wrapf(err, "failed to retreive contracts with project name %s", name)
	}
	networkName := fmt.Sprintf("%s.network", name)

	model.Network = types.Network{
		Name: networkName,
	}

	for _, c := range contracts.NodeContracts {
		contractID, err := strconv.Atoi(c.ContractID)
		if err != nil {
			return types.MachinesModel{}, errors.Wrapf(err, "failed to parse contract with id (%s)", c.ContractID)
		}

		nodeClient, err := client.NcPool.GetNodeClient(client.SubstrateConn, c.NodeID)
		if err != nil {
			return types.MachinesModel{}, errors.Wrapf(err, "failed to get node %d client", c.NodeID)
		}

		dl, err := nodeClient.DeploymentGet(ctx, uint64(contractID))
		if err != nil {
			return types.MachinesModel{}, errors.Wrapf(err, "failed to get deployment with contract id %d", contractID)
		}

		machineMap := map[string]*types.Machine{}
		diskMountPoints := map[string]string{}
		// first get machines and znet
		for idx := range dl.Workloads {
			if dl.Workloads[idx].Type == zos.ZMachineType {
				vm, err := workloads.NewVMFromWorkload(&dl.Workloads[idx], &dl)
				if err != nil {
					return types.MachinesModel{}, errors.Wrapf(err, "failed to parse vm %s data", dl.Workloads[idx].Name)
				}
				machine := machineFromVM(&vm)
				machineMap[machine.Name] = &machine
				for _, mp := range vm.Mounts {
					diskMountPoints[mp.DiskName] = mp.MountPoint
				}
			}
			if dl.Workloads[idx].Type == zos.NetworkType && model.Network.IPRange == "" {
				net, err := workloads.NewNetworkFromWorkload(dl.Workloads[idx], c.NodeID)
				if err != nil {
					return types.MachinesModel{}, errors.Wrapf(err, "failed to parse network %s data", dl.Workloads[idx].Name)
				}

				model.Network.AddWireguardAccess = net.AddWGAccess
				model.Network.IPRange = net.IPRange.String()
			}
		}

		// get disks and qsfss
		for idx := range dl.Workloads {
			if dl.Workloads[idx].Type == zos.ZMountType {
				disk, err := workloads.NewDiskFromWorkload(&dl.Workloads[idx])
				if err != nil {
					return types.MachinesModel{}, errors.Wrapf(err, "failed to parse disk %s data", dl.Workloads[idx].Name)
				}
				machineName, err := getMachineName(disk.Name)
				if err != nil {
					return types.MachinesModel{}, errors.Wrapf(err, "failed to extract machine name from disk with name %s", disk.Name)
				}

				machine, ok := machineMap[machineName]
				if !ok {
					return types.MachinesModel{}, errors.Wrapf(err, "disk (%s) is not mounted on any machine", disk.Name)
				}

				machine.Disks = append(machine.Disks, types.Disk{
					Name:        disk.Name,
					SizeGB:      disk.SizeGB,
					Description: disk.Description,
					MountPoint:  diskMountPoints[disk.Name],
				})
			} else if dl.Workloads[idx].Type == zos.QuantumSafeFSType {
				qsfs, err := workloads.NewQSFSFromWorkload(&dl.Workloads[idx])
				if err != nil {
					return types.MachinesModel{}, errors.Wrapf(err, "failed to parse qsfs %s data", qsfs.Name)
				}

				machineName, err := getMachineName(qsfs.Name)
				if err != nil {
					return types.MachinesModel{}, errors.Wrapf(err, "failed to extract machine name from qsfs with name %s", qsfs.Name)
				}

				machine, ok := machineMap[machineName]
				if !ok {
					return types.MachinesModel{}, errors.Wrapf(err, "qsfs (%s) is not mounted on any machine", qsfs.Name)
				}

				metaBackends := []types.Backend{}
				for _, b := range qsfs.Metadata.Backends {
					metaBackends = append(metaBackends, types.Backend{
						Address:   b.Address,
						Namespace: b.Namespace,
						Password:  b.Password,
					})
				}

				groups := []types.Group{}
				for _, group := range qsfs.Groups {
					bs := types.Backends{}
					for _, b := range group.Backends {
						bs = append(bs, types.Backend{
							Address:   b.Address,
							Namespace: b.Namespace,
							Password:  b.Password,
						})
					}
					groups = append(groups, types.Group{Backends: bs})
				}

				machine.QSFSs = append(machine.QSFSs, types.QSFS{
					MountPoint:           "",
					Description:          qsfs.Description,
					Cache:                qsfs.Cache,
					MinimalShards:        qsfs.MinimalShards,
					ExpectedShards:       qsfs.ExpectedShards,
					RedundantGroups:      qsfs.RedundantGroups,
					RedundantNodes:       qsfs.RedundantNodes,
					MaxZDBDataDirSize:    qsfs.MaxZDBDataDirSize,
					EncryptionAlgorithm:  qsfs.EncryptionAlgorithm,
					EncryptionKey:        qsfs.EncryptionKey,
					CompressionAlgorithm: qsfs.CompressionAlgorithm,
					Metadata: types.Metadata{
						Type:                qsfs.Metadata.Type,
						Prefix:              qsfs.Metadata.Prefix,
						EncryptionAlgorithm: qsfs.Metadata.EncryptionAlgorithm,
						EncryptionKey:       qsfs.Metadata.EncryptionKey,
						Backends:            metaBackends,
					},
					Groups:          groups,
					Name:            qsfs.Name,
					MetricsEndpoint: qsfs.MetricsEndpoint,
				})

			}
		}

		machines := []types.Machine{}
		for _, m := range machineMap {
			machines = append(machines, *m)
		}

		model.Machines = append(model.Machines, machines...)
	}

	model.Name = name

	return model, nil
}

func getMachineName(name string) (string, error) {
	// disk or qsfs name should be in the form: vmname_disk/qsfs_X
	s := strings.Split(name, "_")
	if len(s) == 0 {
		return "", fmt.Errorf("workload name is invalid")
	}
	return s[0], nil
}

func machineFromVM(vm *workloads.VM) types.Machine {
	zlogs := []types.Zlog{}
	for _, zlog := range vm.Zlogs {
		zlogs = append(zlogs, types.Zlog{
			Output: zlog.Output,
		})
	}
	machine := types.Machine{
		NodeID:      0,
		Name:        vm.Name,
		Flist:       vm.Flist,
		PublicIP:    vm.PublicIP,
		PublicIP6:   vm.PublicIP6,
		Planetary:   vm.Planetary,
		Description: vm.Description,
		CPU:         vm.CPU,
		Memory:      vm.Memory,
		RootfsSize:  vm.RootfsSize,
		Entrypoint:  vm.Entrypoint,
		EnvVars:     vm.EnvVars,
		ComputedIP4: vm.ComputedIP,
		ComputedIP6: vm.ComputedIP6,
		WGIP:        vm.IP,
		YggIP:       vm.YggIP,
		Zlogs:       zlogs,
	}
	return machine
}

// func MachineAdd(ctx context.Context, machine types.Machine, projectName string) (types.MachinesModel, error)

// func MachineRemove(ctx context.Context, machineName string, projectName string) (types.MachinesModel, error)
