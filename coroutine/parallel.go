package coroutine

import "sync"

func Parallel(tasks ...func()) {
	var wg sync.WaitGroup
	wg.Add(len(tasks))

	for _, task := range tasks {
		go func(task func()) {
			defer wg.Done()
			task()
		}(task)
	}

	wg.Wait()
}
