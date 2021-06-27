package main

import (
	"github.com/jacoblai/go-coap"
	"runtime"
	"sync"
	"testing"
)

func BenchmarkQps(b *testing.B) {
	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: 12345,
		Payload:   []byte("hello, world!"),
	}

	path := "/a"

	req.SetOption(coap.ETag, "weetag")
	req.SetOption(coap.MaxAge, 3)
	req.SetPathString(path)
	b.ReportAllocs()

	clientPoll := sync.Map{}

	for i := 0; i < runtime.NumCPU()*12; i++ {
		c, err := coap.Dial("udp", "localhost:5683")
		if err != nil {
			b.Error(err)
		}
		c.Busy.Store(false)
		clientPoll.Store(c, nil)
	}

	b.ResetTimer()
	//pres := []byte("hello to you!")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			clientPoll.Range(func(key, value interface{}) bool {
				c := key.(*coap.Conn)
				if c.Busy.Load() == false {
					c.Busy.Store(true)
					_, err := c.Send(req)
					if err != nil {
						c.Busy.Store(false)
						b.Error(err)
					}
					//if bytes.Compare(rv.Payload, pres) != 0 {
					//	c.Busy.Store(false)
					//	b.Error("payload")
					//}
					c.Busy.Store(false)
					return false
				} else {
					return true
				}
			})
		}
	})
}
