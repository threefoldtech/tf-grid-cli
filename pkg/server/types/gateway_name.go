package types

import "github.com/threefoldtech/zos/pkg/gridtypes/zos"

// GatewayNameModel struct for gateway name proxy
type GatewayNameModel struct {
	// Required
	NodeID uint32 `json:"node_id"`
	// Name the fully qualified domain name to use (cannot be present with Name)
	Name string `json:"name"`
	// Backends are list of backend ips
	Backends []zos.Backend `json:"backends"`

	TLSPassthrough bool   `json:"tls_passthrough"`
	Description    string `json:"description"`
	// Optional
	// Passthrough whether to pass tls traffic or not

	// computed

	// FQDN deployed on the node
	// NodeDeploymentID map[uint32]uint64
	FQDN           string `json:"fqdn"`
	NameContractID uint64 `json:"name_contract_id"`
	ContractID     uint64 `json:"contract_id"`
}
