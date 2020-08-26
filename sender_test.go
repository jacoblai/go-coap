package coap

import (
	"testing"
)

func TestSender_Send(t *testing.T) {

	req := Message{
		Type:      Confirmable,
		Code:      GET,
		MessageID: 12345,
		Payload:   []byte("hello, world!"),
	}

	path := "/a"

	req.SetOption(ETag, "weetag")
	req.SetOption(MaxAge, 3)
	req.SetPathString(path)
	s, err := NewSender(8)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 3; i++ {
		res, err := s.Send(req)
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
	}

}
