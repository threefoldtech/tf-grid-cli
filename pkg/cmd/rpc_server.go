package cmd

import (
	"context"

	"github.com/pkg/errors"
	server "github.com/threefoldtech/tf-grid-cli/pkg/server"
)

func TFGridServer() error {
	server, err := server.NewServer()
	if err != nil {
		return errors.Wrap(err, "failed to create new rpc server")
	}

	if err = server.Run(context.Background()); err != nil {
		return errors.Wrap(err, "rpc server stopped")
	}

	return nil
}
