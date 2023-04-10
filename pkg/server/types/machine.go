package types

import (
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

type Machine struct {
	NodeID      uint32            `json:"node_id"`
	FarmID      uint32            `json:"farm_id"`
	Name        string            `json:"name"`
	Flist       string            `json:"flist"`
	PublicIP    bool              `json:"public_ip"`
	PublicIP6   bool              `json:"public_ip6"`
	Planetary   bool              `json:"planetary"`
	Description string            `json:"description"`
	CPU         int               `json:"cpu"`
	Memory      int               `json:"memory"`
	RootfsSize  int               `json:"rootfs_size"`
	Entrypoint  string            `json:"entrypoint"`
	Zlogs       []Zlog            `json:"zlogs"`
	Disks       []Disk            `json:"disks"`
	QSFSs       []QSFS            `json:"qsfss"`
	EnvVars     map[string]string `json:"env_vars"`

	// computed
	ComputedIP4 string `json:"computed_ip4"`
	ComputedIP6 string `json:"computed_ip6"`
	WGIP        string `json:"wireguard_ip"`
	YggIP       string `json:"ygg_ip"`
}

// Zlog logger struct
type Zlog struct {
	Output string `json:"output"`
}

// Disk struct
type Disk struct {
	MountPoint  string `json:"mountpoint"`
	SizeGB      int    `json:"size"`
	Description string `json:"description"`

	// computed
	Name string `json:"name"`
}

// QSFS struct
type QSFS struct {
	MountPoint           string   `json:"mountpoint"`
	Description          string   `json:"description"`
	Cache                int      `json:"cache"`
	MinimalShards        uint32   `json:"minimal_shards"`
	ExpectedShards       uint32   `json:"expected_shards"`
	RedundantGroups      uint32   `json:"redundant_groups"`
	RedundantNodes       uint32   `json:"redundant_nodes"`
	MaxZDBDataDirSize    uint32   `json:"max_zdb_data_dir_size"`
	EncryptionAlgorithm  string   `json:"encryption_algorithm"`
	EncryptionKey        string   `json:"encryption_key"`
	CompressionAlgorithm string   `json:"compression_algorithm"`
	Metadata             Metadata `json:"metadata"`
	Groups               Groups   `json:"groups"`

	// computed
	Name            string `json:"name"`
	MetricsEndpoint string `json:"metrics_endpoint"`
}

// Metadata for QSFS
type Metadata struct {
	Type                string   `json:"type"`
	Prefix              string   `json:"prefix"`
	EncryptionAlgorithm string   `json:"encryption_algorithm"`
	EncryptionKey       string   `json:"encryption_key"`
	Backends            Backends `json:"backends"`
}

// Group is a zos group
type Group struct {
	Backends Backends `json:"backends"`
}

// Backend is a zos backend
type Backend zos.ZdbBackend

// Groups is a list of groups
type Groups []Group

// Backends is a list of backends
type Backends []Backend
