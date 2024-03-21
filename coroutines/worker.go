package coroutines

import "sync"

type Worker struct {
	fns chan func()
	wg  sync.WaitGroup
}

func NewWorker(num int) *Worker {
	s := &Worker{
		fns: make(chan func(), num*2), //nolint:gomnd
	}

	go s.work(num)

	return s
}

func (s *Worker) work(num int) {
	ch := make(chan struct{}, num)
	defer close(ch)

	for {
		select {
		case fn, ok := <-s.fns:
			if !ok {
				return
			}
			ch <- struct{}{}
			go func(fn func()) {
				defer func() {
					<-ch
					s.wg.Done()
				}()

				fn()
			}(fn)
		}
	}
}

func (s *Worker) Push(fns ...func()) {
	s.wg.Add(len(fns))
	for _, fn := range fns {
		s.fns <- fn
	}
}

func (s *Worker) Wait() {
	s.wg.Wait()
}

func (s *Worker) Close() {
	close(s.fns)
}
