package coroutines

type Scheduler struct {
	ch chan struct{}
}

func NewScheduler(max int) *Scheduler {
	return &Scheduler{
		ch: make(chan struct{}, max),
	}
}

func (s *Scheduler) Run(fns ...func()) {
	for _, fn := range fns {
		s.run(fn)
	}
}

func (s *Scheduler) run(fn func()) {
	s.ch <- struct{}{}
	go func() {
		defer func() {
			<-s.ch
		}()
		fn()
	}()
}
