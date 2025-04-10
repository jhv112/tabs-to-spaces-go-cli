package main

import (
	"runtime"
	"sync"
)

// Starts up between 1 and CPU count - 1 consumer coroutines,
// to concurrently consume the producer output.
func produceSyncConsumeAsync[T any](
	produce func(chan<- T) error,
	consume func(<-chan T),
) error {
	ch := make(chan T)

	defer close(ch)

	var wg sync.WaitGroup

	defer wg.Wait()

	consumeAsync := func() {
		defer wg.Done()

		consume(ch)
	}

	// processor count from: https://stackoverflow.com/questions/24073697/a/24073875
	if processorCount := runtime.NumCPU(); processorCount > 1 {
		// do note: many coroutines can be assigned to a single thread; however,
		// each thread will only be running one coroutine simultaneously.
		coroutineCount := processorCount - 1

		wg.Add(coroutineCount)

		for i := 0; i < coroutineCount; i++ {
			go consumeAsync()
		}
	} else {
		wg.Add(1)

		go consumeAsync()
	}

	return produce(ch)
}
