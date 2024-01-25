package coroutine

import "sync"

func Concurrent(limit int, tasks ...func()) {
	var (
		wg sync.WaitGroup
		ch = make(chan struct{}, limit)
	)
	defer close(ch)
	wg.Add(len(tasks))

	for _, task := range tasks {
		ch <- struct{}{}
		go func(task func()) {
			defer func() {
				<-ch
				wg.Done()
			}()
			task()
		}(task)
	}

	wg.Wait()
}
