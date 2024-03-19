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

func (s *Scheduler) Run(fns ...func()) {
	s.wg.Add(len(fns))
	defer s.wg.Wait()

	for _, fn := range fns {
		s.run(fn)
	}
}

func (s *Scheduler) run(fn func()) {
	s.ch <- struct{}{}
	go func() {
		defer func() {
			s.wg.Done()
			<-s.ch
		}()
		fn()
	}()
}
