package mq

import "github.com/pkg/errors"

var (
	KeyNotFound                  = errors.New("MQ configuration key not found")
	KeyAlreadyExists             = errors.New("MQ configuration key already exists")
	GeneralMessageDeliveryFailed = errors.New("general message delivery failed")
	DelayedMessageDeliveryFailed = errors.New("delayed message delivery failed")
	CronMessageDeliveryFailed    = errors.New("cron message delivery failed")
	DelayLevelError              = errors.New("rocketmq delay level error")
)
