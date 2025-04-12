package main

import (
	"runtime"
	"testing"
	"time"
)

func busyLoopFor(d time.Duration) {
	for start := time.Now(); time.Now().Compare(start.Add(d)) < 0; {
	}
}

// No idea, whether this actually works.
// It's supposed to test, whether or not a single processor (preferably a
// single thread) can be used for all coroutines.
func disabledTestRunsOnSingleThread(t *testing.T) {
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()

	defer runtime.UnlockOSThread()

	go func() {
		t.Log("wait fifteen seconds\n")
		busyLoopFor(15 * time.Second)
		t.Log("waited fifteen seconds\n")
	}()

	t.Logf("coroutine count: %d\n", runtime.NumGoroutine())

	t.Log("wait thirty seconds\n")
	busyLoopFor(30 * time.Second)
	t.Log("waited thirty seconds\n")

	t.Fail()
}

func TestConcurrentConsumptionOnSingleThread(t *testing.T) {
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()

	defer runtime.UnlockOSThread()

	observed := make([]int, 0)
	producer := func(ch chan<- int) error {
		for i := 1; i <= 5; i++ {
			ch <- i
		}

		return nil
	}
	consumer := func(ch <-chan int) {
		// no idea, why test doesn't fail as expected, when channel is fed
		// straight into append and defer wg.Wait() is disabled
		i := <-ch

		// better than sleeping for an arbitrary amount of time
		runtime.Gosched()

		observed = append(observed, i)
	}

	produceSyncConsumeAsync(producer, consumer)

	expected := []int{1, 2, 3, 4, 5}

	if len(expected) != len(observed) {
		t.Errorf("want %d elements, have %d elements", len(expected), len(observed))
	}
}
