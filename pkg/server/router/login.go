package router

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/deployer"
)

type Credentials struct {
	Mnemonics string `json:"mnemonics"`
	Network   string `json:"network"`
}

func (r *Router) Login(ctx context.Context, data string) (interface{}, error) {
	cred := Credentials{}

	if err := json.Unmarshal([]byte(data), &cred); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal credentials data")
	}

	if r.Client != nil {
		// TODO: if server already has an initialized client, close old client

	}

	newClient, err := deployer.NewTFPluginClient(cred.Mnemonics, "sr25519", cred.Network, "", "", "", 10, true, false)
	if err != nil {
		return deployer.TFPluginClient{}, errors.Wrap(err, "failed to get tf plugin client")
	}

	r.Client = &newClient
	return nil, nil
}
