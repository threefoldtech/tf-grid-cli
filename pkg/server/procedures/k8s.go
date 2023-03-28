package procedure

import (
	"context"

	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
)

func K8sDeploy(ctx context.Context, cluser types.K8sCluster) (types.K8sCluster, error)

func K8sDelete(ctx context.Context, clusterName string) error

func K8sAddNode(ctx context.Context, cluserName string, node types.K8sNode) (types.K8sCluster, error)

func K8sRemoveNode(ctx context.Context, clusterName string, nodeName string) (types.K8sCluster, error)

func K8sGet(ctx context.Context, clusterName string) (types.K8sCluster, error)
