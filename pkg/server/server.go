package server

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/threefoldtech/grid3-go/deployer"
	router "github.com/threefoldtech/tf-grid-cli/pkg/server/router"
)

type Server struct {
	redisClient RedisClient
	router      Router
	client      *deployer.TFPluginClient
}

type Router struct {
	routes map[string]func(ctx context.Context, client *deployer.TFPluginClient, data string) (interface{}, error)
}

type Request struct {
	JsonRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      string          `json:"id"`
}

// either result or error must has value
type Response struct {
	JsonRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	ID      string      `json:"id"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RedisClient struct {
	Pool *redis.Pool
}

func NewServer() (*Server, error) {
	redisClient, err := NewRedisClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new redis client")
	}

	r := Router{
		routes: make(map[string]func(ctx context.Context, client *deployer.TFPluginClient, data string) (interface{}, error)),
	}

	server := Server{
		redisClient,
		r,
		&deployer.TFPluginClient{},
	}

	server.Register("tfgrid.login", router.Login)

	server.Register("tfgrid.machines.deploy", router.MachinesDeploy)
	server.Register("tfgrid.machines.delete", router.MachinesDelete)
	server.Register("tfgrid.machines.get", router.MachinesGet)
	// server.Register("tfgrid.machines.machine.add", router.MachineAdd)
	// server.Register("tfgrid.machines.machine.remove", router.MachineRemove)

	// server.Register("tfgrid.gateway.name.deploy", router.GatewayNameDeploy)
	// server.Register("tfgrid.gateawy.name.get", router.GatewayNameGet)
	// server.Register("tfgrid.gateway.name.delete", router.GatewayNameDelete)
	// server.Register("tfgrid.gateway.fqdn.deploy", router.GatewayFQDNDeploy)
	// server.Register("tfgrid.gateway.fqdn.get", router.GatewayFQDNGet)
	// server.Register("tfgrid.gateway.fqdn.delete", router.GatewayFQDNDelete)

	server.Register("tfgrid.k8s.get", router.K8sGet)
	server.Register("tfgrid.k8s.deploy", router.K8sDeploy)
	server.Register("tfgrid.k8s.delete", router.K8sDelete)
	// server.Register("tfgrid.k8s.node.add", router.K8sAddNode)
	// server.Register("tfgrid.k8s.node.remove", router.K8sRemoveNode)

	server.Register("tfgrid.zdb.deploy", router.ZDBDeploy)
	server.Register("tfgrid.zdb.delete", router.ZDBDelete)
	server.Register("tfgrid.zdb.get", router.ZDBGet)

	return &server, nil
}

func (s *Server) Register(route string, fn func(context.Context, *deployer.TFPluginClient, string) (interface{}, error)) {
	s.router.routes[route] = fn
}

func NewRedisClient() (RedisClient, error) {
	pool, err := newRedisPool()
	if err != nil {
		return RedisClient{}, errors.Wrap(err, "failed to create new redis pool")
	}
	return RedisClient{
		Pool: pool,
	}, nil
}

func newRedisPool() (*redis.Pool, error) {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) > 10*time.Second {
				_, err := c.Do("PING")
				return err
			}

			return nil
		},
		MaxActive:   100,
		IdleTimeout: 1 * time.Minute,
		Wait:        true,
	}, nil
}

// Run watches a redis queue for incoming messages
func (s *Server) Run(ctx context.Context) error {
	con := s.redisClient.Pool.Get()
	defer con.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			res, err := redis.ByteSlices(con.Do("BRPOP", "tfgrid.client", 0))
			if err != nil {
				return errors.Wrap(err, "failted to read from redis")
			}

			go s.process(ctx, res[1])
		}
	}
}

func (s *Server) process(ctx context.Context, message []byte) {
	args := Request{}
	err := json.Unmarshal(message, &args)
	if err != nil {
		log.Err(err).Msg("failed to unmarshal incoming message. message is dropped")
		return
	}

	err = validateArgs(args)
	if err != nil {
		log.Err(err).Msg("failed to validate incoming message. message is dropped")
		return
	}

	cmd, ok := s.router.routes[args.Method]
	if !ok {
		log.Error().Msgf("invalid command %s. message is dropped", args.Method)
		return
	}

	res, err := cmd(ctx, s.client, string(args.Params))
	response := Response{
		JsonRPC: args.JsonRPC,
		ID:      args.ID,
		Result:  struct{}{},
	}

	if err != nil {
		response.Error = &Error{
			Code:    400,
			Message: err.Error(),
		}
		response.Result = nil
	}

	if res != nil {
		response.Result = res
	}
	con := s.redisClient.Pool.Get()
	defer con.Close()

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Err(err).Msg("failed to marshal response")
		return
	}

	_, err = con.Do("RPUSH", args.ID, responseBytes)
	if err != nil {
		log.Err(err).Msg("failed to push response bytes into redis")
	}
}

func validateArgs(args Request) error {
	// validate jsonrpc standard format
	// any kind of validation on the incoming message should happen here

	// if time.Since(time.Unix(int64(args.Now), 0)) > time.Minute {
	// 	return fmt.Errorf("message with timestamp %d expired", args.Now)
	// }

	return nil
}
