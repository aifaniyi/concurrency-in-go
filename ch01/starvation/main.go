// Package main demonstrates starvation: two goroutines compete for the same
// mutex, but one ("greedy") holds it in one long critical section per loop
// while the other ("polite") repeatedly acquires and releases it in three
// short critical sections per loop. Because Go's mutex doesn't guarantee
// fairness/FIFO ordering, the greedy worker's fewer, longer lock holds win
// the race for the lock far more often, starving the polite worker of CPU
// time and drastically reducing how many loops it completes.
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	wg := &sync.WaitGroup{}
	sharedLock := &sync.Mutex{}
	const runtime = 1 * time.Second

	wg.Add(2)
	go greedyWorker(wg, sharedLock, runtime)
	go politeWorker(wg, sharedLock, runtime)
	wg.Wait()
}

// greedyWorker holds the shared lock for its entire simulated unit of work,
// minimizing how often it has to re-acquire the mutex and re-enter the
// scheduling contention with politeWorker.
func greedyWorker(wg *sync.WaitGroup, sharedLock *sync.Mutex, runtime time.Duration) {
	defer wg.Done()
	var count int

	for begin := time.Now(); time.Since(begin) <= runtime; {
		sharedLock.Lock()
		time.Sleep(3 * time.Nanosecond)
		sharedLock.Unlock()
		count++
	}
	fmt.Printf("Greedy worker was able to execute %v work loops\n", count)
}

// politeWorker splits the same total amount of work into three separate
// lock/unlock cycles per loop iteration. Releasing the lock between each
// small step is "polite" in that it gives other goroutines more chances to
// acquire it, but it also means this worker has to win the lock three times
// as often as greedyWorker to make the same progress - which is exactly why
// it ends up starved.
func politeWorker(wg *sync.WaitGroup, sharedLock *sync.Mutex, runtime time.Duration) {
	defer wg.Done()
	var count int

	for begin := time.Now(); time.Since(begin) <= runtime; {
		sharedLock.Lock()
		time.Sleep(1 * time.Nanosecond)
		sharedLock.Unlock()
		sharedLock.Lock()
		time.Sleep(1 * time.Nanosecond)
		sharedLock.Unlock()
		sharedLock.Lock()
		time.Sleep(1 * time.Nanosecond)
		sharedLock.Unlock()
		count++
	}
	fmt.Printf("Polite worker was able to execute %v work loops.\n", count)
}
