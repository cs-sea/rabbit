package rabbit

import "github.com/streadway/amqp"

type Message struct {
	ID          string
	Body        []byte
	ContentType string
	RoutingKey  string
	Headers     amqp.Table
	Attempt     int
	D           amqp.Delivery
}
