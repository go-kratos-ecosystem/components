package coroutine

import "sync"

type Parallel struct {
	wg sync.WaitGroup
	fs []func()
}

func NewParallel() *Parallel {
	return &Parallel{
		fs: make([]func(), 0),
	}
}

func (p *Parallel) Add(fs ...func()) *Parallel {
	p.fs = append(p.fs, fs...)
	return p
}

func (p *Parallel) Wait() {
	p.wg.Add(len(p.fs))

	for _, f := range p.fs {
		go func(f func()) {
			defer p.wg.Done()
			f()
		}(f)
	}

	p.wg.Wait()
}

func RunParallel(fs ...func()) {
	NewParallel().Add(fs...).Wait()
}
