package server

import (
	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/tf-grid-cli/internal/config"
)

func getClient() (deployer.TFPluginClient, error) {
	cfg, err := config.GetUserConfig()
	if err != nil {
		return deployer.TFPluginClient{}, errors.Wrap(err, "failed to get user configs")
	}

	client, err := deployer.NewTFPluginClient(cfg.Mnemonics, "sr25519", cfg.Network, "", "", "", true, false)
	if err != nil {
		return deployer.TFPluginClient{}, errors.Wrap(err, "failed to get tf plugin client")
	}

	return client, nil
}
