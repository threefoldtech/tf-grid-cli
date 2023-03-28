package procedure

import (
	"github.com/pkg/errors"
	"github.com/threefoldtech/tf-grid-cli/pkg/config"
)

type Credentials struct {
	Mnemonics string `json:"mnemonics"`
	Network   string `json:"network"`
}

func Login(cred Credentials) error {
	path, err := config.GetConfigPath()
	if err != nil {
		return errors.Wrap(err, "failed to get config path")
	}

	cfg := config.Config{}
	cfg.Mnemonics = cred.Mnemonics
	cfg.Network = cred.Network

	err = cfg.Save(path)
	if err != nil {
		return errors.Wrap(err, "failed to save user configs")
	}

	return nil
}
