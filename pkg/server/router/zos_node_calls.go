package router

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	client "github.com/threefoldtech/tf-grid-cli/pkg/server/cli_client"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

func (r *Router) ZOSDeploymentDeploy(ctx context.Context, data string) (interface{}, error) {
	request := client.ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	dl := gridtypes.Deployment{}
	if err := json.Unmarshal([]byte(request.Data), &dl); err != nil {
		return nil, errors.Wrap(err, "failed to parse deployment data")
	}

	return nil, r.client.ZOSDeploymentDeploy(ctx, request.NodeID, dl)
}

func (r *Router) ZOSDeploymentGet(ctx context.Context, data string) (interface{}, error) {
	request := client.ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	contractID := uint64(0)
	if err := json.Unmarshal([]byte(request.Data), &contractID); err != nil {
		return nil, errors.Wrap(err, "failed to parse deployment data")
	}

	return r.client.ZOSDeploymentGet(ctx, request.NodeID, contractID)
}

func (r *Router) ZOSDeploymentDelete(ctx context.Context, data string) (interface{}, error) {
	request := client.ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	contractID := uint64(0)
	if err := json.Unmarshal([]byte(request.Data), &contractID); err != nil {
		return nil, errors.Wrap(err, "failed to parse deployment data")
	}

	return nil, r.client.ZOSDeploymentDelete(ctx, request.NodeID, contractID)
}

func (r *Router) ZOSDeploymentUpdate(ctx context.Context, data string) (interface{}, error) {
	request := client.ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	dl := gridtypes.Deployment{}
	if err := json.Unmarshal([]byte(request.Data), &dl); err != nil {
		return nil, errors.Wrap(err, "failed to parse deployment data")
	}

	return nil, r.client.ZOSDeploymentUpdate(ctx, request.NodeID, dl)
}

func (r *Router) ZOSDeploymentChanges(ctx context.Context, data string) (interface{}, error) {
	request := client.ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	contractID := uint64(0)
	if err := json.Unmarshal([]byte(request.Data), &contractID); err != nil {
		return nil, errors.Wrap(err, "failed to parse deployment data")
	}

	return r.client.ZOSDeploymentChanges(ctx, request.NodeID, contractID)
}

func (r *Router) ZOSStatisticsGet(ctx context.Context, data string) (interface{}, error) {
	request := client.ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	return r.client.ZOSStatisticsGet(ctx, request.NodeID)
}

func (r *Router) ZOSNetworkListWGPorts(ctx context.Context, data string) (interface{}, error) {
	request := client.ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	return r.client.ZOSNetworkListWGPorts(ctx, request.NodeID)
}

func (r *Router) ZOSNetworkInterfaces(ctx context.Context, data string) (interface{}, error) {
	request := client.ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	return r.client.ZOSNetworkInterfaces(ctx, request.NodeID)
}

func (r *Router) ZOSNetworkPublicConfigGet(ctx context.Context, data string) (interface{}, error) {
	request := client.ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	return r.client.ZOSNetworkPublicConfigGet(ctx, request.NodeID)
}

func (r *Router) ZOSSystemDMI(ctx context.Context, data string) (interface{}, error) {
	request := client.ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	return r.client.ZOSSystemDMI(ctx, request.NodeID)
}

func (r *Router) ZOSSystemHypervisor(ctx context.Context, data string) (interface{}, error) {
	request := client.ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	return r.client.ZOSSystemHypervisor(ctx, request.NodeID)
}

func (r *Router) ZOSVersion(ctx context.Context, data string) (interface{}, error) {
	request := client.ZOSNodeRequest{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal deployment data")
	}

	return r.client.ZOSVersion(ctx, request.NodeID)
}
