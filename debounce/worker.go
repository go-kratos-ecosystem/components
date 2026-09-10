package debounce

import (
	"context"
	"time"
)

func (d *Debouncer[T]) run() {
	timer := time.NewTimer(d.delay)
	timer.Stop()
	defer timer.Stop()
	for {
		call, delay, closed := d.next()
		if closed {
			return
		}
		if call != nil {
			d.execute(call)

			continue
		}
		timer.Stop()
		var tick <-chan time.Time
		if delay > 0 {
			timer.Reset(delay)
			tick = timer.C
		}
		select {
		case <-d.wake:
		case <-tick:
		case <-d.ctx.Done():
			d.stop()
		}
	}
}

func (d *Debouncer[T]) next() (*invocation[T], time.Duration, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed || context.Cause(d.ctx) != nil {
		return nil, 0, true
	}
	if d.ready != nil {
		d.running, d.ready = d.ready, nil

		return d.running, 0, false
	}
	if d.pending == nil {
		return nil, 0, false
	}
	if remaining := time.Until(d.pending.due); remaining > 0 {
		return nil, remaining, false
	}
	d.running, d.pending = d.pending, nil

	return d.running, 0, false
}

func (d *Debouncer[T]) execute(call *invocation[T]) {
	err := d.handler(d.ctx, call.value)
	if err != nil && d.config.onError != nil {
		d.config.onError(err)
	}
	d.mu.Lock()
	call.complete(err)
	d.running = nil
	d.mu.Unlock()
}
