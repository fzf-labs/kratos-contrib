package mq

import (
	"context"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
)

type AsynqClient struct {
	log    *log.Helper   //日志
	client *asynq.Client //客户端
}

func NewAsynqClient(
	logger log.Logger,
	redisClientOpt asynq.RedisClientOpt,
) *AsynqClient {
	a := &AsynqClient{
		log:    log.NewHelper(log.With(logger, "module", "mq.asynq.client")),
		client: asynq.NewClient(redisClientOpt),
	}
	return a
}

// ProducerNormalMessage 生产普通消息
func (a *AsynqClient) ProducerNormalMessage(b *MessageConfig, msg []byte) error {
	_, err := a.client.Enqueue(asynq.NewTask(b.Metadata[MetaKeyAsynqQueue], msg))
	if err != nil {
		a.log.Error("Asynq 普通消息推送失败,err:", err)
		return errors.Wrap(GeneralMessageDeliveryFailed, err.Error())
	}
	return nil
}

// ProducerDelayMessage 生产延时消息
func (a *AsynqClient) ProducerDelayMessage(b *MessageConfig, msg []byte, t time.Duration) error {
	_, err := a.client.Enqueue(asynq.NewTask(b.Metadata[MetaKeyAsynqQueue], msg), asynq.ProcessIn(t))
	if err != nil {
		a.log.Error("Asynq 延迟消息推送失败,err:", err)
		return errors.Wrap(DelayedMessageDeliveryFailed, err.Error())
	}
	return nil
}

type AsynqServer struct {
	log             *log.Helper               //日志
	lock            sync.Mutex                //锁
	server          *asynq.Server             //服务端
	scheduler       *asynq.Scheduler          //调度器
	normalConsumers map[*MessageConfig]Handle //普通消费者
	cronConsumers   map[*MessageConfig]string //定时消费者
}

func NwDefaultAsynqConfig() asynq.Config {
	return asynq.Config{
		// 指定要使用多少并发工作人员
		Concurrency: 10,
		// 可以选择指定具有不同优先级的多个队列。
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	}
}

func NewDefaultSchedulerOpts(logger log.Logger) *asynq.SchedulerOpts {
	return &asynq.SchedulerOpts{
		Logger:   log.NewHelper(log.With(logger, "module", "mq.asynq")),
		Location: time.Local, // 使用本地时区
	}
}

// NewAsynqServer
func NewAsynqServer(
	logger log.Logger,
	redisClientOpt asynq.RedisClientOpt,
	asynqConfig asynq.Config,
	schedulerOpts *asynq.SchedulerOpts,
) *AsynqServer {
	a := &AsynqServer{
		log:             log.NewHelper(log.With(logger, "module", "mq.asynq.server")),
		lock:            sync.Mutex{},
		normalConsumers: make(map[*MessageConfig]Handle),
		cronConsumers:   make(map[*MessageConfig]string),
	}
	a.server = asynq.NewServer(
		redisClientOpt,
		asynqConfig,
	)
	a.scheduler = asynq.NewScheduler(redisClientOpt, schedulerOpts)
	return a
}

// ConsumerNormalRegister 注册一个普通消费者
func (a *AsynqServer) ConsumerNormalRegister(b *MessageConfig, handle Handle) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.normalConsumers[b] = handle
}

// ConsumerCronRegister 注册一个定时任务
func (a *AsynqServer) ConsumerCronRegister(b *MessageConfig, handle Handle, cron string) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.normalConsumers[b] = handle
	a.cronConsumers[b] = cron
}

// Start 启动
func (a *AsynqServer) Start(ctx context.Context) error {
	a.log.Info("Asynq server start")
	if len(a.cronConsumers) > 0 {
		for business, cron := range a.cronConsumers {
			b := *business
			_, err := a.scheduler.Register(cron, asynq.NewTask(b.Metadata[MetaKeyAsynqQueue], []byte{}))
			if err != nil {
				a.log.Error("Asynq 定时消息注册失败,err:", err)
				return errors.Wrap(CronMessageDeliveryFailed, err.Error())
			}
		}
		if err := a.scheduler.Start(); err != nil {
			a.log.Error("Asynq 调度器启动失败,err:", err)
			return err
		}
	}
	if len(a.normalConsumers) > 0 {
		mux := asynq.NewServeMux()
		for business, handle := range a.normalConsumers {
			b := *business
			h := handle
			mux.HandleFunc(b.Metadata[MetaKeyAsynqQueue], func(ctx context.Context, task *asynq.Task) error {
				err := h(ctx, task.Payload())
				if err != nil {
					a.log.Error("Asynq 消息业务处理失败,key:", b.Metadata[MetaKeyAsynqQueue], "metadata:", b.Metadata, "body:", string(task.Payload()), "err:", err)
					return err
				}
				return nil
			})
		}
		if err := a.server.Run(mux); err != nil {
			a.log.Error("Asynq服务启动失败,err:", err)
			return err
		}
	}
	return nil
}

func (a *AsynqServer) Stop(ctx context.Context) error {
	a.server.Shutdown()
	a.scheduler.Shutdown()
	a.log.Info("Asynq server stop")
	return nil
}
