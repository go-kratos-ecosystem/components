package coroutine

import "sync"

func Concurrent(limit int, tasks ...func()) {
	var wg sync.WaitGroup
	wg.Add(len(tasks))

	limitCh := make(chan struct{}, limit)

	for _, task := range tasks {
		limitCh <- struct{}{}
		go func(task func()) {
			defer func() {
				<-limitCh
				wg.Done()
			}()
			task()
		}(task)
	}

	wg.Wait()
}
