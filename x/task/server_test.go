package task

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	srv := NewServer()

	go func() {
		wg.Done()
		srv.Start(context.Background())
	}()

	wg.Wait()

	srv.AddTask(func() {
	}, func() {
		fmt.Println("task 2")
	})

	time.Sleep(1 * time.Second)
	srv.Stop(nil)
}

func TestServer_Listen(*testing.T) {
	var wg sync.WaitGroup
	srv := NewServer()
	wg.Add(1)

	go func() {
		wg.Done()
		srv.Start(context.Background())
	}()

	wg.Wait()

	ch := make(chan Task, 1)
	srv.Listen(ch)
	go func() {
		for {
			ch <- func() {
			}
		}
	}()

	time.Sleep(1 * time.Second)
	srv.Stop(nil)
}

func BenchmarkNewServer(b *testing.B) {
	srv := NewServer(
		Size(100),
		Goroutines(1000),
	)
	var counter int64 = 0
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		wg.Done()
		_ = srv.Start(context.Background())
	}()
	wg.Wait()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			srv.AddTask(func() {
				time.Sleep(100 * time.Millisecond)
				atomic.AddInt64(&counter, 1)
			})
		}
	})

	_ = srv.Stop(context.Background())

	assert.Equal(b, int64(b.N), counter)
	b.ReportAllocs()
	log.Println(counter, b.N)
}
