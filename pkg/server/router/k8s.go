package router

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	client "github.com/threefoldtech/tf-grid-cli/pkg/server/cli_client"
)

func (r *Router) K8sDeploy(ctx context.Context, data string) (interface{}, error) {
	cluster := client.K8sCluster{}

	if err := json.Unmarshal([]byte(data), &cluster); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal k8sCluster model data")
	}

	projectName := generateProjectName(cluster.Name)

	cluster, err := r.client.K8sDeploy(ctx, cluster, projectName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deploy cluster")
	}

	return cluster, nil
}

func (r *Router) K8sDelete(ctx context.Context, data string) (interface{}, error) {
	var clusterName string

	if err := json.Unmarshal([]byte(data), &clusterName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal k8sCluster name")
	}

	projectName := generateProjectName(clusterName)

	err := r.client.K8sDelete(ctx, projectName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to delete cluster")
	}

	return nil, nil
}

func (r *Router) K8sGet(ctx context.Context, data string) (interface{}, error) {
	var clusterName string

	if err := json.Unmarshal([]byte(data), &clusterName); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal k8sCluster name")
	}

	projectName := generateProjectName(clusterName)

	cluster, err := r.client.K8sGet(ctx, clusterName, projectName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cluster")
	}

	return cluster, nil
}
