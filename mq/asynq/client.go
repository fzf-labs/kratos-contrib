package asynq

import (
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/hibiken/asynq"
)

type Client struct {
	redisOpt    *asynq.RedisClientOpt
	asynqClient *asynq.Client
}

func NewClient(redisOpt *asynq.RedisClientOpt) *Client {
	c := &Client{
		redisOpt: redisOpt,
	}
	err := c.createAsynqClient()
	if err != nil {
		return nil
	}
	return c
}

func (s *Client) NewTask(typeName string, payload []byte, opts ...asynq.Option) error {
	if s.asynqClient == nil {
		if err := s.createAsynqClient(); err != nil {
			return err
		}
	}

	task := asynq.NewTask(typeName, payload)
	info, err := s.asynqClient.Enqueue(task, opts...)
	if err != nil {
		log.Errorf("[asynq] [%s] Enqueue failed: %s", typeName, err.Error())
		return err
	}
	log.Debugf("[asynq] enqueued task: id=%s queue=%s", info.ID, info.Queue)

	return nil
}

func (s *Client) createAsynqClient() error {
	if s.asynqClient != nil {
		log.Errorf("[asynq] asynq client already created")
		return errors.New("asynq client already created")
	}

	s.asynqClient = asynq.NewClient(s.redisOpt)
	if s.asynqClient == nil {
		log.Errorf("[asynq] create asynq client failed")
		return errors.New("create asynq client failed")
	}

	return nil
}
