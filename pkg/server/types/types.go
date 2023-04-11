package types

type FilterOptions struct {
	FarmID         uint32 `json:"farm_id"`
	PublicConfig   bool   `json:"public_config"`
	PublicIpsCount uint64 `json:"public_ips_count"`
	Dedicated      bool   `json:"dedicated"`
	MRU            uint64 `json:"mru"`
	HRU            uint64 `json:"hru"`
	SRU            uint64 `json:"sru"`
}

type FilterResult struct {
	FilterOption   FilterOptions `json:"filter_options"`
	AvailableNodes []uint32      `json:"available_nodes"`
}
