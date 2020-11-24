package amqp_pool

import (
	"log"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Pool struct {
	connections   map[int]*amqp.Connection
	channels      map[int]*amqp.Channel
	maxConnection int
	maxChannel    int
	busyChannels  map[int]*amqp.Channel
}

func NewPool(maxConnection int, url string) *Pool {
	connections := make(map[int]*amqp.Connection)
	for i := 0; i < maxConnection; i++ {
		c, err := amqp.Dial(url)
		if err != nil {
			log.Fatalf("init connect pool err %v\n", err)
		}
		connections[i] = c
	}
	return &Pool{
		connections: connections,
	}
}

func (p *Pool) GetChannel() (*amqp.Channel, error) {
	for true {
		for _, connection := range p.connections {
			if !connection.IsClosed() {
				return connection.Channel()
			}
		}
	}

	return nil, errors.New("not valid connection")
}

func (p *Pool) ReleaseChannel(channel *amqp.Channel) error {
	return channel.Close()
}
