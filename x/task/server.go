package task

import (
	"context"
	"sync"

	"github.com/go-kratos-ecosystem/components/v2/coroutines"
)

type Task func()

type Server struct {
	tasks  chan Task
	waiter sync.WaitGroup
	stoped chan struct{}

	scheduler *coroutines.Scheduler
}

func NewServer(opts ...Option) *Server {
	o := newOptions(opts...)
	return &Server{
		tasks:     make(chan Task, o.size),
		stoped:    make(chan struct{}),
		scheduler: coroutines.NewScheduler(o.goroutines),
	}
}

func (s *Server) Listen(task <-chan Task) {
	go func() {
		for {
			select {
			case t := <-task:
				s.AddTask(t)
			case <-s.stoped:
				return
			}
		}
	}()
}

func (s *Server) AddTask(tasks ...Task) {
	if s.isStoped() {
		return
	}

	s.waiter.Add(len(tasks))
	for _, task := range tasks {
		s.tasks <- task
	}
}

func (s *Server) Start(context.Context) error {
	for {
		task := <-s.tasks

		s.scheduler.Parallel(func() {
			defer s.waiter.Done()
			task()
		})

		if s.isStoped() && len(s.tasks) == 0 {
			close(s.tasks)
			return nil
		}
	}
}

func (s *Server) isStoped() bool {
	select {
	case <-s.stoped:
		return true
	default:
		return false
	}
}

func (s *Server) Stop(context.Context) error {
	close(s.stoped)
	s.waiter.Wait()
	return nil
}
