package types

import "github.com/threefoldtech/zos/pkg/gridtypes/zos"

// GatewayFQDNModel for gateway FQDN proxy
type GatewayFQDNModel struct {
	// required
	NodeID uint32
	// Backends are list of backend ips
	Backends []zos.Backend
	// FQDN deployed on the node
	FQDN string
	// Name is the workload name
	Name string

	// optional
	// Passthrough whether to pass tls traffic or not
	TLSPassthrough bool
	Description    string

	// SolutionType     string
	// NodeDeploymentID map[uint32]uint64

	// computed
	ContractID uint64
}
