package coroutines

import "sync"

type Scheduler struct {
	ch chan struct{}
	wg sync.WaitGroup
}

func NewScheduler(max int) *Scheduler {
	return &Scheduler{
		ch: make(chan struct{}, max),
	}
}

func (s *Scheduler) Parallel(fns ...func()) {
	for _, fn := range fns {
		s.parallel(fn)
	}
}

func (s *Scheduler) parallel(fn func()) {
	s.ch <- struct{}{}
	go func() {
		defer func() {
			<-s.ch
		}()
		fn()
	}()
}

func (s *Scheduler) Wait(fns ...func()) {
	s.wg.Add(len(fns))
	defer s.wg.Wait()

	for _, fn := range fns {
		s.wait(fn)
	}
}

func (s *Scheduler) wait(fn func()) {
	s.ch <- struct{}{}
	go func() {
		defer func() {
			<-s.ch
			s.wg.Done()
		}()
		fn()
	}()
}
