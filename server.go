// Package coap provides a CoAP client and server.
package coap

import (
	"context"
	"golang.org/x/sys/unix"
	"log"
	"net"
	"runtime"
	"sync"
	"syscall"
	"time"
)

const maxPktLen = 64 * 1024

// Handler is a type that handles CoAP messages.
type Handler interface {
	ServeCOAP(l *net.UDPConn, a *net.UDPAddr, m *Message) *Message
}

var pool = &sync.Pool{
	New: func() interface{} {
		obj := make([]byte, maxPktLen)
		return obj
	},
}

type funcHandler func(l *net.UDPConn, a *net.UDPAddr, m *Message) *Message

func (f funcHandler) ServeCOAP(l *net.UDPConn, a *net.UDPAddr, m *Message) *Message {
	return f(l, a, m)
}

// FuncHandler builds a handler from a function.
func FuncHandler(f func(l *net.UDPConn, a *net.UDPAddr, m *Message) *Message) Handler {
	return funcHandler(f)
}

func handlePacket(l *net.UDPConn, dataLen int, data []byte, u *net.UDPAddr, rh Handler) {
	defer pool.Put(data)
	msg, err := ParseMessage(data[:dataLen])
	if err != nil {
		log.Printf("Error parsing %v", err)
		return
	}

	rv := rh.ServeCOAP(l, u, &msg)
	if rv != nil {
		_ = Transmit(l, u, *rv)
	}
}

// Transmit a message.
func Transmit(l *net.UDPConn, a *net.UDPAddr, m Message) error {
	d, err := m.MarshalBinary()
	if err != nil {
		return err
	}

	if a == nil {
		_, err = l.Write(d)
	} else {
		_, err = l.WriteTo(d, a)
	}

	return err
}

// Receive a message.
func Receive(l *net.UDPConn, buf []byte) (Message, error) {
	l.SetReadDeadline(time.Now().Add(ResponseTimeout))

	nr, _, err := l.ReadFromUDP(buf)
	if err != nil {
		return Message{}, err
	}
	return ParseMessage(buf[:nr])
}

// ListenAndServe binds to the given address and serve requests forever.
func ListenAndServe(n, addr string, rh Handler) error {
	_, err := net.ResolveUDPAddr(n, addr)
	if err != nil {
		return err
	}

	for i := 0; i < runtime.NumCPU()*2; i++ {
		l, err := newUDPListener(addr)
		if err != nil {
			return err
		}
		err = Serve(l, rh)
		if err != nil {
			return err
		}
	}

	return nil
}

func newUDPListener(address string) (conn *net.UDPConn, err error) {
	cfgFn := func(network, address string, conn syscall.RawConn) (err error) {
		// 构造fd控制函数
		fn := func(fd uintptr) {
			// 设置REUSEPORT
			err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, unix.SO_REUSEPORT, 1)
			if err != nil {
				return
			}
		}

		if err = conn.Control(fn); err != nil {
			return
		}
		return
	}

	lc := net.ListenConfig{Control: cfgFn}
	lp, err := lc.ListenPacket(context.Background(), "udp", address)
	if err != nil {
		return
	}
	conn = lp.(*net.UDPConn)
	return
}

// Serve processes incoming UDP packets on the given listener, and processes
// these requests forever (or until the listener is closed).
func Serve(listener *net.UDPConn, rh Handler) error {
	buf := pool.Get().([]byte)
	for {
		nr, addr, err := listener.ReadFromUDP(buf)
		if err != nil {
			if neterr, ok := err.(net.Error); ok && (neterr.Temporary() || neterr.Timeout()) {
				time.Sleep(5 * time.Millisecond)
				continue
			}
			return err
		}
		go handlePacket(listener, nr, buf, addr, rh)
	}
}
