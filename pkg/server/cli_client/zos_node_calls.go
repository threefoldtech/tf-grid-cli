package client

import (
	"context"
	"net"

	"github.com/pkg/errors"
	client "github.com/threefoldtech/grid3-go/node"
	"github.com/threefoldtech/zos/pkg/capacity/dmi"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

type ZOSNodeRequest struct {
	NodeID uint32 `json:"node_id"`
	Data   string `json:"data"`
}

type Statistics struct {
	Total gridtypes.Capacity `json:"total"`
	Used  gridtypes.Capacity `json:"used"`
}

func (r *CLIClient) ZOSDeploymentDeploy(ctx context.Context, nodeID uint32, dl gridtypes.Deployment) error {

	nodeClient, err := r.client.GetNodeClient(nodeID)
	if err != nil {
		return errors.Wrapf(err, "failed to get node %d client", nodeID)
	}

	return nodeClient.DeploymentDeploy(ctx, dl)
}

func (r *CLIClient) ZOSDeploymentGet(ctx context.Context, nodeID uint32, contractID uint64) (gridtypes.Deployment, error) {

	nodeClient, err := r.client.GetNodeClient(nodeID)
	if err != nil {
		return gridtypes.Deployment{}, errors.Wrapf(err, "failed to get node %d client", nodeID)
	}

	deployment, err := nodeClient.DeploymentGet(ctx, contractID)
	if err != nil {
		return gridtypes.Deployment{}, errors.Wrapf(err, "failed to get deployment with contract id %d", contractID)
	}

	return deployment, nil
}

func (r *CLIClient) ZOSDeploymentDelete(ctx context.Context, nodeID uint32, contractID uint64) error {

	nodeClient, err := r.client.GetNodeClient(nodeID)
	if err != nil {
		return errors.Wrapf(err, "failed to get node %d client", nodeID)
	}

	err = nodeClient.DeploymentDelete(ctx, contractID)
	if err != nil {
		return errors.Wrapf(err, "failed to delete deployment with contract id %d", contractID)
	}

	return nil
}

func (r *CLIClient) ZOSDeploymentUpdate(ctx context.Context, nodeID uint32, dl gridtypes.Deployment) error {

	nodeClient, err := r.client.GetNodeClient(nodeID)
	if err != nil {
		return errors.Wrapf(err, "failed to get node %d client", nodeID)
	}

	err = nodeClient.DeploymentUpdate(ctx, dl)
	if err != nil {
		return errors.Wrapf(err, "failed to update deployment with contract id %d", dl.ContractID)
	}

	return nil
}

func (r *CLIClient) ZOSDeploymentChanges(ctx context.Context, nodeID uint32, contractID uint64) ([]gridtypes.Workload, error) {
	nodeClient, err := r.client.GetNodeClient(nodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %d client", nodeID)
	}

	workloads, err := nodeClient.DeploymentChanges(ctx, contractID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get changes for deployment with contract id %d", contractID)
	}

	return workloads, nil
}

func (r *CLIClient) ZOSStatisticsGet(ctx context.Context, nodeID uint32) (Statistics, error) {
	nodeClient, err := r.client.GetNodeClient(nodeID)
	if err != nil {
		return Statistics{}, errors.Wrapf(err, "failed to get node %d client", nodeID)
	}

	total, used, err := nodeClient.Statistics(ctx)
	if err != nil {
		return Statistics{}, errors.Wrapf(err, "failed to get statistics for node with id %d", nodeID)
	}

	return Statistics{
		Total: total,
		Used:  used,
	}, nil
}

func (r *CLIClient) ZOSNetworkListWGPorts(ctx context.Context, nodeID uint32) ([]uint16, error) {
	nodeClient, err := r.client.GetNodeClient(nodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %d client", nodeID)
	}

	ports, err := nodeClient.NetworkListWGPorts(ctx)
	if err != nil {
		return nil, err
	}

	return ports, nil
}

func (r *CLIClient) ZOSNetworkInterfaces(ctx context.Context, nodeID uint32) (map[string][]net.IP, error) {
	nodeClient, err := r.client.GetNodeClient(nodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %d client", nodeID)
	}

	ips, err := nodeClient.NetworkListInterfaces(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get statistics for node with id %d", nodeID)
	}

	return ips, nil
}

func (r *CLIClient) ZOSNetworkPublicConfigGet(ctx context.Context, nodeID uint32) (client.PublicConfig, error) {
	nodeClient, err := r.client.GetNodeClient(nodeID)
	if err != nil {
		return client.PublicConfig{}, errors.Wrapf(err, "failed to get node %d client", nodeID)
	}

	cfg, err := nodeClient.NetworkGetPublicConfig(ctx)
	if err != nil {
		return client.PublicConfig{}, errors.Wrapf(err, "failed to get statistics for node with id %d", nodeID)
	}

	return cfg, nil
}

func (r *CLIClient) ZOSSystemDMI(ctx context.Context, nodeID uint32) (dmi.DMI, error) {
	nodeClient, err := r.client.GetNodeClient(nodeID)
	if err != nil {
		return dmi.DMI{}, errors.Wrapf(err, "failed to get node %d client", nodeID)
	}

	resDMI, err := nodeClient.SystemDMI(ctx)
	if err != nil {
		return dmi.DMI{}, errors.Wrapf(err, "failed to get statistics for node with id %d", nodeID)
	}

	return resDMI, nil
}

func (r *CLIClient) ZOSSystemHypervisor(ctx context.Context, nodeID uint32) (string, error) {
	nodeClient, err := r.client.GetNodeClient(nodeID)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get node %d client", nodeID)
	}

	res, err := nodeClient.SystemHypervisor(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get statistics for node with id %d", nodeID)
	}

	return res, nil
}

func (r *CLIClient) ZOSVersion(ctx context.Context, nodeID uint32) (client.Version, error) {
	nodeClient, err := r.client.GetNodeClient(nodeID)
	if err != nil {
		return client.Version{}, errors.Wrapf(err, "failed to get node %d client", nodeID)
	}

	version, err := nodeClient.SystemVersion(ctx)
	if err != nil {
		return client.Version{}, errors.Wrapf(err, "failed to get statistics for node with id %d", nodeID)
	}

	return version, nil
}
