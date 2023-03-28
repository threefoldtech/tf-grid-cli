package types

// ZDB workload struct
type ZDB struct {
	NodeID      uint32
	Name        string
	Password    string
	Public      bool
	Size        int
	Description string
	Mode        string
	Port        uint32
	Namespace   string

	// computed
	IPs []string
}
