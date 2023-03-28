package procedure

// type NetworkParams struct {
// 	Name               string `json:"name"`
// 	IPRange            string `json:"ip_range"`
// 	AddWireguardAccess bool   `json:"add_wireguard_access"`
// 	Description        string `json:"description"`
// }

// type MachinesModel struct {
// 	Name        string        `json:"name"`
// 	Network     NetworkParams `json:"network"`
// 	Machines    []Machine     `json:"machines"`
// 	Metadata    string        `json:"metadata"`
// 	Description string        `json:"description"`
// }

// type Machine struct {
// 	Name       string `json:"name"`
// 	NodeID     uint32 `json:"node_id"`
// 	Disks      []Disk `json:"disks"`
// 	PublicIP   bool   `json:"public_ip"`
// 	Planetary  bool   `json:"planetary"`
// 	CPU        uint32 `json:"cpu"`
// 	Memory     uint64 `json:"memory"`
// 	RootFSSize uint64 `json:"rootfs_size"`
// 	Flist      string `json:"flist"`
// 	Entrypoint string `json:"entrypoint"`
// 	SSHKey     string `json:"ssh_key"`
// }

// type Disk struct {
// 	Name       string `json:"name"`
// 	Size       uint32 `json:"size"`
// 	Mountpoint string `json:"mountpoint"`
// }

// type MachinesResult struct {
// 	NetworkResult NetworkResult   `json:"network_result"`
// 	MachineResult []MachineResult `json:"machine_result"`
// }

// type NetworkResult struct {
// 	WireguardConfig string `json:"wireguard_config"`
// }

// type MachineResult struct {
// 	Name      string `json:"name"`
// 	PublicIP  string `json:"public_ip"`
// 	PublicIP6 string `json:"public_ip6"`
// 	YggIP     string `json:"ygg_ip"`
// }

// func MachinesDeploy(ctx context.Context, model MachinesModel, client deployer.TFPluginClient) (MachinesResult, error) {

// 	vms, disks, network, err := extractWorkloads(model)
// 	if err != nil {
// 		return MachinesResult{}, errors.Wrap(err, "failed to extract workloads data")
// 	}

// 	resVMs, wgConfgig, err := client.DeployMachines(ctx, model.Name, vms, disks, network)
// 	if err != nil {
// 		return MachinesResult{}, errors.Wrap(err, "failed to deploy machines")
// 	}

// 	machinesResult := getMachinesResult(resVMs, wgConfgig)
// 	return machinesResult, nil
// }

// func extractWorkloads(m MachinesModel) ([]workloads.VM, []workloads.Disk, workloads.ZNet, error) {
// 	vms := []workloads.VM{}
// 	disks := []workloads.Disk{}
// 	for _, vm := range m.Machines {
// 		mounts := []workloads.Mount{}
// 		for _, disk := range vm.Disks {
// 			disks = append(disks, workloads.Disk{
// 				Name:   disk.Name,
// 				SizeGB: int(disk.Size),
// 			})
// 			mounts = append(mounts, workloads.Mount{
// 				DiskName:   disk.Name,
// 				MountPoint: disk.Mountpoint,
// 			})
// 		}
// 		vms = append(vms, workloads.VM{
// 			Name:       vm.Name,
// 			Flist:      vm.Flist,
// 			PublicIP:   vm.PublicIP,
// 			Planetary:  vm.Planetary,
// 			CPU:        int(vm.CPU),
// 			Memory:     int(vm.Memory),
// 			RootfsSize: int(vm.RootFSSize),
// 			Entrypoint: vm.Entrypoint,
// 			EnvVars: map[string]string{
// 				"SSH_KEY": vm.SSHKey,
// 			},
// 			NetworkName: m.Network.Name,
// 			Mounts:      mounts,
// 		})
// 	}
// 	ip, err := gridtypes.ParseIPNet(m.Network.IPRange)
// 	if err != nil {
// 		return nil, nil, workloads.ZNet{}, errors.Wrap(err, "failed to parse ip range")
// 	}
// 	network := workloads.ZNet{
// 		Name:        m.Network.Name,
// 		Description: m.Network.Description,
// 		IPRange:     ip,
// 		AddWGAccess: m.Network.AddWireguardAccess,
// 	}
// 	return vms, disks, network, nil
// }

// func getMachinesResult(vms []workloads.VM, wgConfig string) MachinesResult {
// 	machinesResult := MachinesResult{}
// 	for _, vm := range vms {
// 		machinesResult.MachineResult = append(machinesResult.MachineResult, MachineResult{
// 			Name:      vm.Name,
// 			PublicIP:  vm.ComputedIP,
// 			PublicIP6: vm.ComputedIP6,
// 			YggIP:     vm.YggIP,
// 		})
// 	}
// 	machinesResult.NetworkResult.WireguardConfig = wgConfig
// 	return machinesResult
// }
