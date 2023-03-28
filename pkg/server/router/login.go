package router

import (
	"context"
)

func Login(ctx context.Context, data string) (string, error) {
	// cred := procedure.Credentials{}

	// if err := json.Unmarshal([]byte(data), &cred); err != nil {
	// 	return server.Response{
	// 		Err: errors.Wrap(err, "failed to unmarshal credentials data"),
	// 	}
	// }

	// if err := procedure.Login(cred); err != nil {
	// 	return server.Response{
	// 		Err: errors.Wrap(err, "failed to login"),
	// 	}
	// }

	// return server.Response{}
	return "", nil
}
