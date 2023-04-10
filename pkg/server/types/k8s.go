package types

import "github.com/threefoldtech/zos/pkg/gridtypes"

// K8sCluster struct for k8s cluster
type K8sCluster struct {
	Name        string    `json:"name"`
	Master      *K8sNode  `json:"master"`
	Workers     []K8sNode `json:"workers"`
	Token       string    `json:"token"`
	NetworkName string    `json:"network_name"`
	SSHKey      string    `json:"ssh_key"`

	//optional

	//computed
	NodeDeploymentID map[uint32]uint64          `json:"node_deployment_id"`
	NodesIPRange     map[uint32]gridtypes.IPNet `json:"nodes_ip_range"`
}

// K8sNode kubernetes data
type K8sNode struct {
	Name      string `json:"name"`
	NodeID    uint32 `json:"node_id"`
	FarmID    uint32 `json:"farm_id"`
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
