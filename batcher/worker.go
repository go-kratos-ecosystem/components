package batcher

import (
	"context"
	"time"
)

type worker[T any] struct {
	ctx     context.Context
	handler Handler[T]
	config  config

	batch    []T
	timer    *time.Timer
	tick     <-chan time.Time
	firstErr error
}

func (w *worker[T]) run(requests <-chan request[T]) error {
	w.timer = time.NewTimer(w.config.flushInterval)
	w.timer.Stop()
	defer w.timer.Stop()

	for {
		select {
		case req, ok := <-requests:
			if !ok {
				w.flush()

				return w.firstErr
			}
			w.accept(req)
		case <-w.tick:
			w.flush()
		}
	}
}

func (w *worker[T]) accept(req request[T]) {
	if req.flush != nil {
		w.flush()
		req.flush.complete(w.firstErr)

		return
	}
	if len(w.batch) == 0 {
		w.batch = make([]T, 0, w.config.batchSize)
		w.timer.Reset(w.config.flushInterval)
		w.tick = w.timer.C
	}
	w.batch = append(w.batch, req.item)
	if len(w.batch) == w.config.batchSize {
		w.flush()
	}
}

func (w *worker[T]) flush() {
	if len(w.batch) == 0 {
		return
	}
	w.timer.Stop()
	w.tick = nil
	batch := w.batch
	w.batch = nil
	err := w.handler(w.ctx, batch)
	if err != nil && w.firstErr == nil {
		w.firstErr = err
	}
	if w.config.onResult != nil {
		w.config.onResult(Result{Size: len(batch), Err: err})
	}
}
