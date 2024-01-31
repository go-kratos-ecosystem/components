package coroutine

import "sync"

func Wait(fs ...func()) {
	var wg sync.WaitGroup
	wg.Add(len(fs))
	for _, f := range fs {
		go func(f func()) {
			defer wg.Done()
			f()
		}(f)
	}
	wg.Wait()
}
