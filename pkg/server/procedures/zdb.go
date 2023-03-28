package procedure

import (
	"context"

	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
)

func ZDBDeploy(ctx context.Context, zdb types.ZDB) (types.ZDB, error)

func ZDBDelete(ctx context.Context, name string) error

func ZDBGet(ctx context.Context, name string) (types.ZDB, error)
