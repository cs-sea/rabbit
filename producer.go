package rabbit

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type Producer interface {
	Publish(ctx context.Context, exchange *Exchange, message *Message) error
}

type ProducerService struct {
	Rabbit RabbitService
}

func NewProducerService(rabbit RabbitService) Producer {
	return &ProducerService{Rabbit: rabbit}
}

func (p *ProducerService) Publish(ctx context.Context, exchange *Exchange, message *Message) error {
	// conn pool
	connPool := p.Rabbit.GetConnectionPool()
	conn, err := connPool.Get()
	if err != nil {
		return err
	}
	defer connPool.Release(conn)

	// channel pool
	pool := conn.GetChannelPool()
	ch, err := pool.Get()
	if err != nil {
		return err
	}
	defer pool.Release(ch)

	if err = p.Rabbit.ExchangeDeclare(ch, exchange); err != nil {
		time.Sleep(time.Minute)
		log.Errorln(err)
		return err
	}

	return ch.Publish(&PublishData{
		Exchange:  exchange.Name,
		Key:       message.RoutingKey,
		Mandatory: false,
		Immediate: false,
		Msg: &amqp.Publishing{
			Headers:      message.Headers,
			ContentType:  message.ContentType,
			MessageId:    message.ID,
			DeliveryMode: amqp.Persistent,
			Body:         message.Body,
		},
	})
}
