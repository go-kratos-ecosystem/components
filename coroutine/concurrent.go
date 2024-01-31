package coroutine

import "sync"

type Concurrent struct {
	wg sync.WaitGroup
	ch chan struct{}
	fs []func()
}

func NewConcurrent(limit int) *Concurrent {
	return &Concurrent{
		ch: make(chan struct{}, limit),
		fs: make([]func(), 0),
	}
}

func (c *Concurrent) Add(fs ...func()) *Concurrent {
	c.fs = append(c.fs, fs...)
	return c
}

func (c *Concurrent) Wait() {
	c.wg.Add(len(c.fs))

	for _, f := range c.fs {
		c.ch <- struct{}{}
		go func(f func()) {
			defer func() {
				<-c.ch
				c.wg.Done()
			}()
			f()
		}(f)
	}

	c.wg.Wait()
}

func (c *Concurrent) Close() {
	close(c.ch)
}

func RunConcurrent(limit int, tasks ...func()) {
	c := NewConcurrent(limit).Add(tasks...)
	defer c.Close()
	c.Wait()
}
