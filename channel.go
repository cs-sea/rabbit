package amqp_pool

import (
	"github.com/streadway/amqp"
)

type Channel struct {
	ch *amqp.Channel
}

type PublishData struct {
	Exchange  string
	Key       string
	Mandatory bool // 如果没有队列 true 返回给生产者 false 会丢掉
	Immediate bool // 如果没有消费者 true 返回给生产者，false 丢入队列
	Msg       amqp.Publishing
}

// name, kind string, durable, autoDelete, internal, noWait bool, args Table
type ExchangeDeclareConfig struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

// name string, durable, autoDelete, exclusive, noWait bool, args Table
type QueueDeclareConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type ChannelService interface {
	Publish() error
}

func NewChannel(ch *amqp.Channel) *Channel {
	return &Channel{ch: ch}
}

func (c *Channel) Publish(data *PublishData) error {
	return c.ch.Publish(data.Exchange, data.Key, data.Mandatory, data.Immediate, data.Msg)
}

func (c *Channel) ExchangeDeclare(config *ExchangeDeclareConfig) error {
	return c.ch.ExchangeDeclare(config.Name, config.Kind, config.Durable, config.AutoDelete, config.Internal, config.NoWait, config.Args)
}

func (c *Channel) QueueDeclare(config *QueueDeclareConfig) (amqp.Queue, error) {
	return c.ch.QueueDeclare(config.Name, config.Durable, config.AutoDelete, config.Exclusive, config.NoWait, config.Args)
}
