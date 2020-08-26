package coap

import (
	"errors"
	"log"
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
	if r.current < r.qmax {
		atomic.AddInt64(&r.current, 1)
	} else {
		atomic.StoreInt64(&r.current, 0)
	}
	log.Println(r.current)
	c := r.clients[r.current-1]
	return c.Send(req)
}
