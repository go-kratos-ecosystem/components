package coroutines

import "sync"

func Wait(fns ...func()) {
	var wg sync.WaitGroup
	wg.Add(len(fns))
	for _, fn := range fns {
		go func(fn func()) {
			defer wg.Done()
			fn()
		}(fn)
	}
	wg.Wait()
}

func Parallel(max int, fns ...func()) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, max)
	defer close(ch)
	wg.Add(len(fns))
	for _, fn := range fns {
		ch <- struct{}{}
		go func(task func()) {
			defer func() {
				<-ch
				wg.Done()
			}()
			task()
		}(fn)
	}
	wg.Wait()
}
