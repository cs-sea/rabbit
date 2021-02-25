package rabbit

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
	Msg       *amqp.Publishing
}

// name, kind string, durable, autoDelete, internal, noWait bool, args Table
type Exchange struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

// name string, durable, autoDelete, exclusive, noWait bool, args Table
type Queue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type ChannelService interface {
	Publish(data *PublishData) error
}

func NewChannel(ch *amqp.Channel) *Channel {
	return &Channel{ch: ch}
}

func (c *Channel) Publish(data *PublishData) error {
	return c.ch.Publish(data.Exchange, data.Key, data.Mandatory, data.Immediate, *data.Msg)
}

func (c *Channel) ExchangeDeclare(ex *Exchange) error {
	return c.ch.ExchangeDeclare(ex.Name, ex.Kind, ex.Durable, ex.AutoDelete, ex.Internal, ex.NoWait, ex.Args)
}

func (c *Channel) QueueDeclare(queue *Queue) (amqp.Queue, error) {
	return c.ch.QueueDeclare(queue.Name, queue.Durable, queue.AutoDelete, queue.Exclusive, queue.NoWait, queue.Args)
}

func (c *Channel) QueueBind(queueName, key, exchangeName string, args amqp.Table) error {
	return c.ch.QueueBind(queueName, key, exchangeName, true, args)
}

func (c *Channel) Consume(name string) (<-chan amqp.Delivery, error) {
	return c.ch.Consume(name, "", false, false, false, false, nil)
}
