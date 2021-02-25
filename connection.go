package rabbit

import (
	"container/list"
	"github.com/streadway/amqp"
	"log"
)

type Connection struct {
	conn        *amqp.Connection
	channelPool *ChannelPool
}

type ConnectionConfig struct {
	ChanNum int
	Url     string
}

func NewConnection(config *ConnectionConfig) *Connection {
	c, err := amqp.Dial(config.Url)

	if err != nil {
		log.Fatalln(err)
	}

	chList := list.New()
	for i := 0; i < config.ChanNum; i++ {
		ch, err := c.Channel()
		if err != nil {
			log.Fatalln(err)
		}

		mch := NewChannel(ch)
		chList.PushBack(mch)
	}

	pool := NewChannelPool(chList)

	return &Connection{
		conn:        c,
		channelPool: pool,
	}
}

func (c *Connection) GetChannelPool() *ChannelPool {
	return c.channelPool
}
