package main

import (
	"bytes"
	"github.com/jacoblai/go-coap"
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

	c, err := coap.Dial("udp", "localhost:5683")
	if err != nil {
		b.Error(err)
	}
	pres := []byte("hello to you!")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rv, err := c.Send(req)
			if err != nil {
				b.Error(err)
			}

			if bytes.Compare(rv.Payload, pres) != 0 {
				b.Error("payload")
			}
		}
	})
}
