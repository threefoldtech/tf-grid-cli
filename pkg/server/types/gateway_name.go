package types

import "github.com/threefoldtech/zos/pkg/gridtypes/zos"

// GatewayNameModel struct for gateway name proxy
type GatewayNameModel struct {
	// Required
	NodeID uint32
	// Name the fully qualified domain name to use (cannot be present with Name)
	Name string
	// Backends are list of backend ips
	Backends []zos.Backend

	TLSPassthrough bool
	Description    string
	// Optional
	// Passthrough whether to pass tls traffic or not

	// computed

	// FQDN deployed on the node
	// NodeDeploymentID map[uint32]uint64
	FQDN           string
	NameContractID uint64
	ContractID     uint64
}
