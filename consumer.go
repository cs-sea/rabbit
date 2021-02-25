package rabbit

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
)

type Consumer interface {
	Consume(name string) error
}

type ConsumerService struct {
	Rabbit RabbitService
	ch     *Channel
}

func NewConsumerService(rabbit RabbitService) *ConsumerService {
	return &ConsumerService{Rabbit: rabbit}
}

func (c *ConsumerService) Consume(name string) error {
	// 连接池，不过消费者好像没啥作用。。。
	connPool := c.Rabbit.GetConnectionPool()
	conn, err := connPool.Get()

	if err != nil {
		return err
	}
	defer connPool.Release(conn)

	// channel 池
	pool := conn.GetChannelPool()
	ch, err := pool.Get()

	if err != nil {
		return err
	}
	defer pool.Release(ch)

	deliveries, err := ch.Consume(name)
	if err != nil {
		log.Fatalln(err)
	}

	go c.handle(deliveries, ch)

	return nil
}

func (c *ConsumerService) Listen(name string) {
	deliveries, err := c.ch.Consume(name)
	if err != nil {
		log.Fatalln(err)
	}

	go c.handle(deliveries, c.ch)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
}

func (c *ConsumerService) handle(deliveries <-chan amqp.Delivery, ch *Channel) {
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)

		attempt, ok := d.Headers["rabbit_attempt"]
		if !ok {
			attempt = 1
		}

		log.Println(attempt)
		if d.Headers == nil {
			d.Headers = make(map[string]interface{})
		}
		d.Headers[RabbitAttempt] = attempt
		d.Ack(false)
		c.release(ch, &Message{
			Body:       d.Body,
			RoutingKey: d.RoutingKey,
			Headers:    d.Headers,
			Attempt:    0,
		})
	}

	log.Printf("handle: deliveries channel closed")
}

func (c *ConsumerService) release(ch *Channel, message *Message) {
	ch.Publish(&PublishData{
		Exchange:  "",
		Key:       message.RoutingKey,
		Mandatory: false,
		Immediate: false,
		Msg: &amqp.Publishing{
			Headers:      message.Headers,
			DeliveryMode: amqp.Persistent,
			Body:         message.Body,
		},
	})
}
