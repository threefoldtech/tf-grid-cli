package types

import "github.com/threefoldtech/zos/pkg/gridtypes/zos"

type Machine struct {
	NodeID      uint32
	Name        string
	Flist       string
	PublicIP    bool
	PublicIP6   bool
	Planetary   bool
	Description string
	CPU         int
	Memory      int
	RootfsSize  int
	Entrypoint  string
	Zlogs       []Zlog
	Disks       []Disk
	QSFSs       []QSFS
	EnvVars     map[string]string

	// computed
	ComputedIP4 string
	ComputedIP6 string
	WGIP        string
	YggIP       string
}

// Zlog logger struct
type Zlog struct {
	Output string
}

// Disk struct
type Disk struct {
	Name        string
	MountPoint  string
	SizeGB      int
	Description string
}

// QSFS struct
type QSFS struct {
	Name                 string
	MountPoint           string
	Description          string
	Cache                int
	MinimalShards        uint32
	ExpectedShards       uint32
	RedundantGroups      uint32
	RedundantNodes       uint32
	MaxZDBDataDirSize    uint32
	EncryptionAlgorithm  string
	EncryptionKey        string
	CompressionAlgorithm string
	Metadata             Metadata
	Groups               Groups

	// computed
	MetricsEndpoint string
}

// Metadata for QSFS
type Metadata struct {
	Type                string
	Prefix              string
	EncryptionAlgorithm string
	EncryptionKey       string
	Backends            Backends
}

// Group is a zos group
type Group struct {
	Backends Backends
}

// Backend is a zos backend
type Backend zos.ZdbBackend

// Groups is a list of groups
type Groups []Group

// Backends is a list of backends
type Backends []Backend

// type MachineResult struct {
// 	NodeID      uint32
// 	Name        string
// 	Flist       string
// 	PublicIP    bool
// 	PublicIP6   bool
// 	Planetary   bool
// 	Description string
// 	CPU         int
// 	Memory      int
// 	RootfsSize  int
// 	Entrypoint  string
// 	Zlogs       []Zlog
// 	Disks       []Disk
// 	QSFSs       []QSFS
// 	EnvVars     map[string]string

// 	// computed
// 	ComputedIP4 string
// 	ComputedIP6 string
// 	WGIP        string
// 	YggIP       string
// }
