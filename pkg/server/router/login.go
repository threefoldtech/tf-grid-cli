package router

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	procedure "github.com/threefoldtech/tf-grid-cli/pkg/server/procedures"
)

func Login(ctx context.Context, data string) (interface{}, error) {
	cred := procedure.Credentials{}

	if err := json.Unmarshal([]byte(data), &cred); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal credentials data")
	}

	if err := procedure.Login(cred); err != nil {
		return nil, errors.Wrap(err, "failed to login")
	}

	return struct{}{}, nil
}
