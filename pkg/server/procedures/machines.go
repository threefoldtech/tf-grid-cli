package procedure

import (
	"context"

	"github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
)

// machines.deploy
func MachinesDeploy(ctx context.Context, model types.MachinesModel, client deployer.TFPluginClient) (types.MachinesModel, error)

// machines.deelte
func MachinesDelete(ctx context.Context, name string) error

func MachineAdd(ctx context.Context, machine types.Machine, projectName string) (types.MachinesModel, error)

func MachineRemove(ctx context.Context, machineName string, projectName string) (types.MachinesModel, error)

func MachinesGet(ctx context.Context, name string) (types.MachinesModel, error)
