package server

import (
	"context"
)

type Router struct {
	routes map[string]func(ctx context.Context, data string) (string, error)
}

// func Login(ctx context.Context, data string) Response {
// 	cred := procedure.Credentials{}

// 	if err := json.Unmarshal([]byte(data), &cred); err != nil {
// 		return Response{
// 			Err: errors.Wrap(err, "failed to unmarshal credentials data"),
// 		}
// 	}

// 	if err := procedure.Login(cred); err != nil {
// 		return Response{
// 			Err: errors.Wrap(err, "failed to login"),
// 		}
// 	}

// 	return Response{}
// }

// func MachinesDeploy(ctx context.Context, data string) Response {
// 	machinesModel := procedure.MachinesModel{}
// 	if err := json.Unmarshal([]byte(data), &machinesModel); err != nil {
// 		return Response{
// 			Err: errors.Wrap(err, "failed to unmarshal machines data"),
// 		}
// 	}

// 	client, err := getClient()
// 	if err != nil {
// 		return Response{
// 			Err: errors.Wrap(err, "failed to get tf grid client"),
// 		}
// 	}

// 	res, err := procedure.MachinesDeploy(ctx, machinesModel, client)
// 	if err != nil {
// 		return Response{
// 			Err: err,
// 		}
// 	}

// 	resultBytes, err := json.Marshal(res)
// 	if err != nil {
// 		return Response{
// 			Err: errors.Wrap(err, "failed to marshal machines result"),
// 		}
// 	}

// 	return Response{
// 		Result: string(resultBytes),
// 	}
// }

// func MachinesGet(ctx context.Context, data string) Response {
// 	projectName := data

// 	client, err := getClient()
// 	if err != nil {
// 		return Response{
// 			Err: err,
// 		}
// 	}

// 	res, err := procedure.MachinesGet(ctx, projectName, client)
// 	if err != nil {
// 		return Response{
// 			Err: err,
// 		}
// 	}

// 	resBytes, err := json.Marshal(res)
// 	if err != nil {
// 		return Response{
// 			Err: errors.Wrap(err, "failed to marshal deployments"),
// 		}
// 	}

// 	return Response{
// 		Result: string(resBytes),
// 	}
// }

// func MachinesDelete(ctx context.Context, data string) Response {
// 	projectName := data

// 	client, err := getClient()
// 	if err != nil {
// 		return Response{
// 			Err: err,
// 		}
// 	}

// 	if err := procedure.MachinesDelete(ctx, projectName, client); err != nil {
// 		return Response{
// 			Err: err,
// 		}
// 	}

// 	return Response{}
// }
