package amqp_pool

import (
	"container/list"
	"sync"
)

type ConnectionPool struct {
	m        *sync.Cond
	connList *list.List
}
type ConnPoolConfig struct {
	ConnNum uint16
	Url     string
}

func NewConnectionPool(config *ConnPoolConfig) *ConnectionPool {
	m := sync.NewCond(new(sync.Mutex))

	connList := list.New()

	for i := uint16(0); i < config.ConnNum; i++ {

		conn := NewConnection(&ConnectionConfig{
			ChanNum: 5,
			Url:     "amqp://guest:guest@localhost:5672/",
		})

		connList.PushBack(conn)
	}

	return &ConnectionPool{m: m, connList: connList}
}

func (c *ConnectionPool) Get() (*Connection, error) {

	c.m.L.Lock()
	defer c.m.L.Unlock()

	for {
		conn := c.connList.Front()

		if conn != nil {
			c.connList.Remove(conn)
			return conn.Value.(*Connection), nil
		}

		c.m.Wait()
	}
}

func (c *ConnectionPool) Release(conn *Connection) {
	c.m.L.Lock()
	defer c.m.L.Unlock()

	c.connList.PushBack(conn)
	c.m.Signal()
}
