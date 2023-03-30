package router

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/deployer"
	procedure "github.com/threefoldtech/tf-grid-cli/pkg/server/procedures"
)

func Login(ctx context.Context, client *deployer.TFPluginClient, data string) (interface{}, error) {
	cred := procedure.Credentials{}

	if err := json.Unmarshal([]byte(data), &cred); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal credentials data")
	}

	if client != nil {
		// TODO: if server already has an initialized client, close old client, assign new one

	}

	newClient, err := deployer.NewTFPluginClient(cred.Mnemonics, "sr25519", cred.Network, "", "", "", true, false)
	if err != nil {
		return deployer.TFPluginClient{}, errors.Wrap(err, "failed to get tf plugin client")
	}

	*client = newClient
	return struct{}{}, nil
}
