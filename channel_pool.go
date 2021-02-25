package rabbit

import (
	"container/list"
	"sync"
)

type ChannelPool struct {
	pool *list.List

	m *sync.Cond
}

func NewChannelPool(pool *list.List) *ChannelPool {
	l := new(sync.Mutex)
	m := sync.NewCond(l)
	return &ChannelPool{
		m:    m,
		pool: pool,
	}
}

func (p *ChannelPool) Get() (*Channel, error) {

	p.m.L.Lock()
	defer p.m.L.Unlock()
	for {
		e := p.pool.Front()
		if e != nil {
			p.pool.Remove(e)
			return e.Value.(*Channel), nil
		}

		p.m.Wait()

	}
}

func (p *ChannelPool) Release(channel *Channel) {

	p.m.L.Lock()
	defer p.m.L.Unlock()
	p.pool.PushBack(channel)
	p.m.Signal()
}
