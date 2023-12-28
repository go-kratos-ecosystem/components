package udp

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	var (
		server *Server
		wg     sync.WaitGroup
		done   = make(chan []byte, 1)
	)

	wg.Add(2)

	go func() {
		defer wg.Done()

		server = NewServer(":12190", WithHandler(func(conn net.PacketConn, buf []byte, addr net.Addr) {
			done <- buf
		}), WithRecoveryHandler(func(conn net.PacketConn, buf []byte, addr net.Addr, err interface{}) {
			t.Log(err)
		}), WithBufSize(1024))

		go server.Start(context.Background())

		time.Sleep(time.Second * 5)
		server.Stop(context.Background())
	}()

	go func() {
		defer wg.Done()

		c, err := net.Dial("udp", ":12190")
		if err != nil {
			t.Error(err)
			return
		}
		defer c.Close()

		_, err = c.Write([]byte("test"))
		if err != nil {
			t.Error(err)
			return
		}
	}()

	wg.Wait()

	buf := <-done
	if string(buf) != "test" {
		t.Fatal("buf not equal test")
	}
}
