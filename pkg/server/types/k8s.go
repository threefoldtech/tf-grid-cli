package types

import "github.com/threefoldtech/zos/pkg/gridtypes"

// K8sCluster struct for k8s cluster
type K8sCluster struct {
	Name        string
	Master      *K8sNode
	Workers     []K8sNode
	Token       string
	NetworkName string
	SSHKey      string

	//optional

	//computed
	NodeDeploymentID map[uint32]uint64
	NodesIPRange     map[uint32]gridtypes.IPNet
}

// K8sNode kubernetes data
type K8sNode struct {
	Name      string
	NodeID    uint32
	DiskSize  int
	PublicIP  bool
	PublicIP6 bool
	Planetary bool
	Flist     string
	CPU       int
	Memory    int

	// computed
	ComputedIP4 string
	ComputedIP6 string
	WGIP        string
	YggIP       string
}
