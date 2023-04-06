package router

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
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

func (r *Router) ZOSDeploymentDeploy(ctx context.Context, data string) (interface{}, error) {
	request := ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	nodeClient, err := r.client.GetNodeClient(request.NodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %d client", request.NodeID)
	}

	dl := gridtypes.Deployment{}
	if err := json.Unmarshal([]byte(request.Data), &dl); err != nil {
		return nil, errors.Wrap(err, "failed to parse deployment data")
	}

	return nil, nodeClient.DeploymentDeploy(ctx, dl)
}

func (r *Router) ZOSDeploymentGet(ctx context.Context, data string) (interface{}, error) {
	request := ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	nodeClient, err := r.client.GetNodeClient(request.NodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %d client", request.NodeID)
	}

	contractID := uint64(0)
	if err := json.Unmarshal([]byte(request.Data), &contractID); err != nil {
		return nil, errors.Wrap(err, "failed to parse deployment data")
	}

	deployment, err := nodeClient.DeploymentGet(ctx, contractID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get deployment with contract id %d", contractID)
	}

	return deployment, nil
}

func (r *Router) ZOSDeploymentDelete(ctx context.Context, data string) (interface{}, error) {
	request := ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	nodeClient, err := r.client.GetNodeClient(request.NodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %d client", request.NodeID)
	}

	contractID := uint64(0)
	if err := json.Unmarshal([]byte(request.Data), &contractID); err != nil {
		return nil, errors.Wrap(err, "failed to parse deployment data")
	}

	err = nodeClient.DeploymentDelete(ctx, contractID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to delete deployment with contract id %d", contractID)
	}

	return nil, nil
}

func (r *Router) ZOSDeploymentUpdate(ctx context.Context, data string) (interface{}, error) {
	request := ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	nodeClient, err := r.client.GetNodeClient(request.NodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %d client", request.NodeID)
	}

	dl := gridtypes.Deployment{}
	if err := json.Unmarshal([]byte(request.Data), &dl); err != nil {
		return nil, errors.Wrap(err, "failed to parse deployment data")
	}

	err = nodeClient.DeploymentUpdate(ctx, dl)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update deployment with contract id %d", dl.ContractID)
	}

	return nil, nil
}

func (r *Router) ZOSDeploymentChanges(ctx context.Context, data string) (interface{}, error) {
	request := ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	nodeClient, err := r.client.GetNodeClient(request.NodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %d client", request.NodeID)
	}

	contractID := uint64(0)
	if err := json.Unmarshal([]byte(request.Data), &contractID); err != nil {
		return nil, errors.Wrap(err, "failed to parse deployment data")
	}

	workloads, err := nodeClient.DeploymentChanges(ctx, contractID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get changes for deployment with contract id %d", contractID)
	}

	return workloads, nil
}

func (r *Router) ZOSStatisticsGet(ctx context.Context, data string) (interface{}, error) {
	request := ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	nodeClient, err := r.client.GetNodeClient(request.NodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %d client", request.NodeID)
	}

	total, used, err := nodeClient.Statistics(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get statistics for node with id %d", request.NodeID)
	}

	return Statistics{
		Total: total,
		Used:  used,
	}, nil
}

func (r *Router) ZOSNetworkListWGPorts(ctx context.Context, data string) (interface{}, error) {
	request := ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	nodeClient, err := r.client.GetNodeClient(request.NodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %d client", request.NodeID)
	}

	ports, err := nodeClient.NetworkListWGPorts(ctx)
	if err != nil {
		return nil, err
	}

	return ports, nil
}

func (r *Router) ZOSNetworkInterfaces(ctx context.Context, data string) (interface{}, error) {
	request := ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	nodeClient, err := r.client.GetNodeClient(request.NodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %d client", request.NodeID)
	}

	ips, err := nodeClient.NetworkListInterfaces(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get statistics for node with id %d", request.NodeID)
	}

	return ips, nil
}

func (r *Router) ZOSNetworkPublicConfigGet(ctx context.Context, data string) (interface{}, error) {
	request := ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	nodeClient, err := r.client.GetNodeClient(request.NodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %d client", request.NodeID)
	}

	cfg, err := nodeClient.NetworkGetPublicConfig(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get statistics for node with id %d", request.NodeID)
	}

	return cfg, nil
}

func (r *Router) ZOSSystemDMI(ctx context.Context, data string) (interface{}, error) {
	request := ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	nodeClient, err := r.client.GetNodeClient(request.NodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %d client", request.NodeID)
	}

	dmi, err := nodeClient.SystemDMI(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get statistics for node with id %d", request.NodeID)
	}

	return dmi, nil
}

func (r *Router) ZOSSystemHypervisor(ctx context.Context, data string) (interface{}, error) {
	request := ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	nodeClient, err := r.client.GetNodeClient(request.NodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %d client", request.NodeID)
	}

	res, err := nodeClient.SystemHypervisor(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get statistics for node with id %d", request.NodeID)
	}

	return res, nil
}

func (r *Router) ZOSVersion(ctx context.Context, data string) (interface{}, error) {
	request := ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	nodeClient, err := r.client.GetNodeClient(request.NodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %d client", request.NodeID)
	}

	version, err := nodeClient.SystemVersion(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get statistics for node with id %d", request.NodeID)
	}

	return version, nil
}
