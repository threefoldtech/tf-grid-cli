package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Server struct {
	redisClient RedisClient
	router      Router
}

type Request struct {
	RetQueue string `json:"ret_queue"`
	Now      uint64 `json:"now"`
	Cmd      string `json:"cmd"`
	Data     string `json:"data"`
}

type Response struct {
	Result string `json:"result"`
	Err    error  `json:"error"`
}

type RedisClient struct {
	Pool *redis.Pool
}

func NewServer() (*Server, error) {
	redisClient, err := NewRedisClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new redis client")
	}

	router := Router{
		routes: make(map[string]func(ctx context.Context, data string) (response Response)),
	}

	server := Server{
		redisClient,
		router,
	}

	server.Register("login", Login)
	server.Register("machines.deploy", MachinesDeploy)
	server.Register("machines.get", MachinesGet)
	server.Register("machines.delete", MachinesDelete)

	return &server, nil
}

func (s *Server) Register(route string, fn func(context.Context, string) Response) {
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

	cmd, ok := s.router.routes[args.Cmd]
	if !ok {
		log.Error().Msg("invalid command. message is dropped")
		return
	}

	resopnse := cmd(ctx, args.Data)
	b, err := json.Marshal(resopnse)
	if err != nil {
		log.Err(err).Msg("failed to marshal response")
		return
	}

	con := s.redisClient.Pool.Get()
	defer con.Close()

	_, err = con.Do("RPUSH", args.RetQueue, b)
	if err != nil {
		log.Err(err).Msg("failed to push response bytes into redis")
	}
}

func validateArgs(args Request) error {
	// any kind of validation on the incoming message should happen here

	if time.Since(time.Unix(int64(args.Now), 0)) > time.Minute {
		return fmt.Errorf("message with timestamp %d expired", args.Now)
	}

	return nil
}
