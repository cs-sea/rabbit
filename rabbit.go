package rabbit

import (
	"sync"
)

type RabbitService interface {
	ExchangeDeclare(ch *Channel, exchange *Exchange) error
	QueueDeclare(ch *Channel, queue *Queue) error
	BindQueue(ch *Channel, exchangeName, key, queueName string) error
	GetConnectionPool() *ConnectionPool
}

const RabbitAttempt = "rabbit_attempt"

type Rabbit struct {
	connections *ConnectionPool

	exchanges *sync.Map
	queues    *sync.Map
}

type Options struct {
	MaxConnection int
	Host          string
	MaxChannels   int
}

func NewRabbit(opts *Options) *Rabbit {
	connections := NewConnectionPool(&ConnPoolOptions{
		ConnNum: opts.MaxConnection,
		Url:     opts.Host,
	})

	exchanges := new(sync.Map)
	queues := new(sync.Map)

	return &Rabbit{connections: connections, exchanges: exchanges, queues: queues}
}

// declare exchange if not exist
func (r *Rabbit) ExchangeDeclare(ch *Channel, exchange *Exchange) error {
	// exist will return
	if _, ok := r.exchanges.Load(exchange.Name); ok {
		return nil
	}

	if err := ch.ExchangeDeclare(exchange); err != nil {
		return err
	}

	r.exchanges.Store(exchange.Name, exchange.Name)

	return nil
}

// declare queue if not exist
func (r *Rabbit) QueueDeclare(ch *Channel, queue *Queue) error {
	if _, ok := r.queues.Load(queue.Name); ok {
		return nil
	}

	if _, err := ch.QueueDeclare(queue); err != nil {
		return err
	}

	r.queues.Store(queue.Name, queue.Name)
	return nil
}

// bind exchange with queue
func (r *Rabbit) BindQueue(ch *Channel, exchangeName, key, queueName string) error {
	return ch.QueueBind(queueName, key, exchangeName, nil)
}

func (r *Rabbit) GetConnectionPool() *ConnectionPool {
	return r.connections
}
