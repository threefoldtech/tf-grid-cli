package router

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	procedure "github.com/threefoldtech/tf-grid-cli/pkg/server/procedures"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
)

func K8sDeploy(ctx context.Context, data string) (interface{}, error) {
	cluster := types.K8sCluster{}

	if err := json.Unmarshal([]byte(data), &cluster); err != nil {
		return types.K8sCluster{}, errors.Wrap(err, "failed to unmarshal k8sCluster model data")
	}

	client, err := getClient()
	if err != nil {
		return types.K8sCluster{}, errors.Wrap(err, "failed to create new grid client")
	}

	cluster, err = procedure.K8sDeploy(ctx, cluster, &client)
	if err != nil {
		return types.K8sCluster{}, errors.Wrap(err, "failed to deploy cluster")
	}

	return cluster, nil
}

func K8sDelete(ctx context.Context, data string) (interface{}, error) {
	var clusterName string

	if err := json.Unmarshal([]byte(data), &clusterName); err != nil {
		return struct{}{}, errors.Wrap(err, "failed to unmarshal k8sCluster name")
	}

	client, err := getClient()
	if err != nil {
		return struct{}{}, errors.Wrap(err, "failed to create new grid client")
	}

	err = procedure.K8sDelete(ctx, clusterName, &client)
	if err != nil {
		return struct{}{}, errors.Wrap(err, "failed to delete cluster")
	}

	return struct{}{}, nil
}

// func K8sAddNode(ctx context.Context, data string) (string, error)

// // func K8sRemoveNode(ctx context.Context, data string) (string, error)

func K8sGet(ctx context.Context, data string) (interface{}, error) {
	cluster := types.K8sCluster{}
	var clusterName string

	if err := json.Unmarshal([]byte(data), &clusterName); err != nil {
		return types.K8sCluster{}, errors.Wrap(err, "failed to unmarshal k8sCluster name")
	}

	client, err := getClient()
	if err != nil {
		return types.K8sCluster{}, errors.Wrap(err, "failed to create new grid client")
	}

	cluster, err = procedure.K8sGet(ctx, clusterName, &client)
	if err != nil {
		return types.K8sCluster{}, errors.Wrap(err, "failed to get cluster")
	}

	return cluster, nil
}
