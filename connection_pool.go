package rabbit

import (
	"container/list"
	"sync"
)

type ConnectionPool struct {
	m        *sync.Cond
	connList *list.List
}
type ConnPoolOptions struct {
	ConnNum    int
	Url        string
	ChannelNum int
}

const (
	DefaultConnNum    = 1
	DefaultChannelNum = 1

	DefaultHost = "amqp://guest:guest@localhost:5672/"
)

func NewConnectionPool(opts *ConnPoolOptions) *ConnectionPool {
	setDefaultConfig(opts)

	m := sync.NewCond(new(sync.Mutex))
	connList := list.New()

	for i := 0; i < opts.ConnNum; i++ {

		conn := NewConnection(&ConnectionConfig{
			ChanNum: opts.ChannelNum,
			Url:     opts.Url,
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

func setDefaultConfig(opt *ConnPoolOptions) {
	if opt.ConnNum == 0 {
		opt.ConnNum = DefaultConnNum
	}

	if opt.ChannelNum == 0 {
		opt.ChannelNum = DefaultChannelNum
	}

	if opt.Url == "" {
		opt.Url = DefaultHost
	}
}
