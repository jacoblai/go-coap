package coap

import (
	"errors"
	"sync/atomic"
)

type Sender struct {
	current int64
	qmax    int64
	clients []*Conn
}

func NewSender(q int64) (*Sender, error) {
	if q <= 0 {
		return nil, errors.New("q error")
	}
	re := &Sender{
		current: 0,
		qmax:    q,
		clients: make([]*Conn, 0),
	}
	for i := 0; i < int(q); i++ {
		c, err := Dial("udp", "localhost:5683")
		if err != nil {
			return nil, err
		}
		re.clients = append(re.clients, c)
	}
	return re, nil
}

func (r *Sender) Send(req Message) (*Message, error) {
	current := atomic.LoadInt64(&r.current)
	max := atomic.LoadInt64(&r.qmax)
	if current == max {
		atomic.StoreInt64(&current, 0)
	}
	c := r.clients[current]
	atomic.AddInt64(&current, 1)
	return c.Send(req)
}
