package router

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	client "github.com/threefoldtech/tf-grid-cli/pkg/server/cli_client"
)

func (r *Router) Login(ctx context.Context, data string) (interface{}, error) {
	cred := client.Credentials{}

	if err := json.Unmarshal([]byte(data), &cred); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal credentials data")
	}

	if err := r.client.Login(ctx, cred); err != nil {
		return nil, errors.Wrap(err, "failed to login")
	}

	return nil, nil
}
