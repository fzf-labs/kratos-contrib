package mq

import (
	"context"
	"time"
)

// Handle 消费者业务方法
type Handle func(ctx context.Context, msg []byte) error

type Client interface {
	// ProducerNormalMessage 生产普通消息
	ProducerNormalMessage(b *MessageConfig, msg []byte) error
	// ProducerDelayMessage 生产延时消息
	ProducerDelayMessage(b *MessageConfig, msg []byte, t time.Duration) error
}

type Server interface {
	// ConsumerNormalRegister 注册一个普通消费者
	ConsumerNormalRegister(b *MessageConfig, handle Handle)
	// ConsumerCronRegister 注册一个定时任务
	ConsumerCronRegister(b *MessageConfig, handle Handle, cron string)
	// Start 启动
	Start(context.Context) error
	// Stop 停止
	Stop(context.Context) error
}
