package types

// Machines model ensures that each node has one deployment that includes all workloads
type MachinesModel struct {
	Name        string    `json:"name"`     // this is the project name, should be unique
	Network     Network   `json:"network"`  // network specs
	Machines    []Machine `json:"machines"` // machines specs
	Metadata    string    `json:"metadata"`
	Description string    `json:"description"`

	// computed
	NodeDeploymentID map[uint32]uint64
}

type Network struct {
	AddWireguardAccess bool   `json:"add_wireguard_access"` // true to add access node
	IPRange            string `json:"ip_range"`

	// computed
	Name            string `json:"name"` // network name will be (projectname.network)
	WireguardConfig string
}
