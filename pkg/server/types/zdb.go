package types

type ZDB struct {
	NodeID      uint32 `json:"node_id"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Public      bool   `json:"public"`
	Size        int    `json:"size"`
	Description string `json:"description"`
	Mode        string `json:"mode"`
	Port        uint32 `json:"port"`
	Namespace   string `json:"namespace"`

	// computed
	IPs []string `json:"ips"`
}
