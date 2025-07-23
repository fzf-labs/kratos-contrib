package mq

import (
	"sync"
)

// MQType 消息队列类型
type MQType string

const (
	MQTypeAsynq    MQType = "asynq"    // asynq
	MQTypeRocketMQ MQType = "rocketmq" // rocketmq
	MQTypeKafka    MQType = "kafka"    // kafka
	MQTypeRabbitMQ MQType = "rabbitmq" // rabbitmq
)

// 元数据键定义
type MetaKey string

const (
	MetaKeyAsynqQueue         MetaKey = "asynq_queue"          // asynq 队列名称
	MetaKeyRocketMQTag        MetaKey = "rocketmq_tag"         // rocketmq 标签
	MetaKeyRocketMQGroupId    MetaKey = "rocketmq_group_id"    // rocketmq 组 ID
	MetaKeyKafkaGroupId       MetaKey = "kafka_group_id"       // kafka 组 ID
	MetaKeyKafkaPartition     MetaKey = "kafka_partition"      // kafka 分区
	MetaKeyRabbitMQExchange   MetaKey = "rabbitmq_exchange"    // rabbitmq 交换机
	MetaKeyRabbitMQRoutingKey MetaKey = "rabbitmq_routing_key" // rabbitmq 路由键
	MetaKeyRabbitMQQueue      MetaKey = "rabbitmq_queue"       // rabbitmq 队列
)

// MessageConfig 消息配置结构体
type MessageConfig struct {
	Key      string             `json:"key"`
	Metadata map[MetaKey]string `json:"metadata"`
}

// MessageConfigManager 配置管理器
type MessageConfigManager struct {
	mu      sync.RWMutex
	Type    MQType                    `json:"type"`
	Configs map[string]*MessageConfig `json:"configs"`
}

// NewMessageConfigManager 创建配置管理器
func NewMessageConfigManager(mqType MQType) *MessageConfigManager {
	return &MessageConfigManager{
		Type:    mqType,
		Configs: make(map[string]*MessageConfig),
	}
}

// Register 注册配置
func (c *MessageConfigManager) Register(config *MessageConfig) (*MessageConfig, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.Configs[config.Key]; exists {
		return nil, KeyAlreadyExists
	}
	c.Configs[config.Key] = config
	return config, nil
}

// GetMessageConfig 获取消息配置
func (c *MessageConfigManager) GetMessageConfig(key string) (*MessageConfig, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	config, exists := c.Configs[key]
	if !exists {
		return nil, KeyNotFound
	}
	return config, nil
}
