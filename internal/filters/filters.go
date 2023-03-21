package filters

import (
	"fmt"

	"github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/grid_proxy_server/pkg/client"
	"github.com/threefoldtech/grid_proxy_server/pkg/types"
)

func GetAvailableNode(client client.Client, filter types.NodeFilter) (uint32, error) {
	nodes, err := deployer.FilterNodes(client, filter)
	if err != nil {
		return 0, err
	}
	if len(nodes) == 0 {
		return 0, fmt.Errorf(
			"no node with free resources available using node filter: farmIDs: %v, mru: %d, sru: %d, freeips: %d, domain: %t",
			filter.FarmIDs,
			*filter.FreeMRU,
			*filter.FreeSRU,
			*filter.FreeIPs,
			*filter.Domain,
		)
	}

	node := uint32(nodes[0].NodeID)
	return node, nil
}

func BuildK8sFilter(k8sNode workloads.K8sNode, farmID uint64, k8sNodesNum uint) types.NodeFilter {
	freeMRUs := uint64(k8sNode.Memory*int(k8sNodesNum)) / 1024
	freeSRUs := uint64(k8sNode.DiskSize * int(k8sNodesNum))
	freeIPs := uint64(0)
	if k8sNode.PublicIP {
		freeIPs = uint64(k8sNodesNum)
	}

	return buildGenericFilter(freeMRUs, freeSRUs, freeIPs, farmID, false)
}

func BuildVMFilter(vm workloads.VM, disk workloads.Disk, farmID uint64) types.NodeFilter {
	freeMRUs := uint64(vm.Memory) / 1024
	freeSRUs := uint64(vm.RootfsSize) / 1024
	freeIPs := uint64(0)
	if vm.PublicIP {
		freeIPs = 1
	}
	freeSRUs += uint64(disk.SizeGB)
	return buildGenericFilter(freeMRUs, freeSRUs, freeIPs, farmID, false)
}

func BuildGatewayFilter(farmID uint64) types.NodeFilter {
	return buildGenericFilter(0, 0, 0, farmID, true)
}

func buildGenericFilter(mrus, srus, ips, farmID uint64, domain bool) types.NodeFilter {
	status := "up"
	return types.NodeFilter{
		Status:  &status,
		FreeMRU: &mrus,
		FreeSRU: &srus,
		FreeIPs: &ips,
		FarmIDs: []uint64{farmID},
		Domain:  &domain,
	}
}
