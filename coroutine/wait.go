package coroutine

import "sync"

func Wait(fn func()) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		fn()
		wg.Done()
	}()
	wg.Wait()
}
