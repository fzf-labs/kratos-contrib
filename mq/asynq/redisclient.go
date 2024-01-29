package asynq

import (
	"crypto/tls"
	"time"

	"github.com/hibiken/asynq"
)

const (
	defaultRedisAddress = "127.0.0.1:6379"
)

func NewRedisClientOpt(opts ...RedisClientOption) *asynq.RedisClientOpt {
	r := &asynq.RedisClientOpt{
		Addr: defaultRedisAddress,
		DB:   0,
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(r)
		}
	}
	return r
}

type RedisClientOption func(o *asynq.RedisClientOpt)

func WithAddress(addr string) RedisClientOption {
	return func(s *asynq.RedisClientOpt) {
		s.Addr = addr
	}
}

func WithRedisAuth(userName, password string) RedisClientOption {
	return func(s *asynq.RedisClientOpt) {
		s.Username = userName
		s.Password = password
	}
}

func WithRedisPoolSize(size int) RedisClientOption {
	return func(s *asynq.RedisClientOpt) {
		s.PoolSize = size
	}
}

func WithRedisDatabase(db int) RedisClientOption {
	return func(s *asynq.RedisClientOpt) {
		s.DB = db
	}
}

func WithDialTimeout(timeout time.Duration) RedisClientOption {
	return func(s *asynq.RedisClientOpt) {
		s.DialTimeout = timeout
	}
}

func WithReadTimeout(timeout time.Duration) RedisClientOption {
	return func(s *asynq.RedisClientOpt) {
		s.ReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) RedisClientOption {
	return func(s *asynq.RedisClientOpt) {
		s.WriteTimeout = timeout
	}
}

func WithTLSConfig(c *tls.Config) RedisClientOption {
	return func(s *asynq.RedisClientOpt) {
		s.TLSConfig = c
	}
}
