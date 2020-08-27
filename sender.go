package coap

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
)

type Sender struct {
	current int64
	qmax    int64
	clients sync.Map
}

func NewSender(q int64) (*Sender, error) {
	if q <= 0 {
		return nil, errors.New("q error")
	}
	re := &Sender{
		current: 0,
		qmax:    q,
	}
	for i := 0; i < int(q); i++ {
		c, err := Dial("udp", "localhost:5683")
		if err != nil {
			return nil, err
		}
		re.clients.Store(i, c)
	}
	return re, nil
}

func (r *Sender) Send(req Message) (*Message, error) {
	current := atomic.LoadInt64(&r.current)
	max := atomic.LoadInt64(&r.qmax)
	if current == max {
		atomic.StoreInt64(&current, 0)
	}
	cc, ok := r.clients.Load(current)
	if !ok {
		return nil, errors.New("pool err")
	}
	atomic.AddInt64(&current, 1)
	c := cc.(*Conn)
	log.Println(current)
	return c.Send(req)
}
