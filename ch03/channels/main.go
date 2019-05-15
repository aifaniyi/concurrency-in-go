// Package main demonstrates a graceful-shutdown producer/consumer pipeline
// built on channels: a single producer emits jobs on a timer, a pool of
// worker goroutines consumes them, and an OS signal (Ctrl+C / SIGTERM)
// triggers the producer to close the channel so every worker's range loop
// exits cleanly and the program can shut down without leaking goroutines.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	workerNumber := 5

	wg := &sync.WaitGroup{}
	wg.Add(workerNumber)

	// Buffered channel (capacity 1) so a signal sent before main starts
	// listening isn't lost, as recommended by the signal package docs.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Buffered so the producer can stay ahead of the workers by up to
	// workerNumber jobs without blocking on every send.
	jobsStream := make(chan *job, workerNumber)

	// create consumers
	for i := 0; i < workerNumber; i++ {
		go worker(wg, i, jobsStream)
	}

	// create produder
	go producer(jobsStream, sigs)

	// Blocks until every worker has returned from its range loop, which
	// only happens once the producer closes jobsStream on shutdown.
	wg.Wait()
	fmt.Println("stopped all go routines")
}

// job is the unit of work handed from the producer to the workers.
type job struct {
	id int
}

// worker consumes jobs from jobStream until the channel is closed and
// drained, then returns. Ranging over a channel is the idiomatic way to
// consume until closure: the loop exits automatically once the channel is
// both closed and empty.
func worker(wg *sync.WaitGroup, id int, jobStream <-chan *job) {
	defer wg.Done()

	fmt.Printf("starting worker %d\n", id)

	for job := range jobStream {
		time.Sleep(5 * time.Millisecond)
		fmt.Printf("worker %d processed job %d\n", id, job.id)
	}

	fmt.Printf("stopping worker %d\n", id)
}

// producer emits a new job every tick until it receives a termination
// signal, at which point it closes jobStream. Closing on the producer side
// (never the consumer side) is the safe convention for a single-writer
// channel: it guarantees no further sends will happen after close, so
// workers ranging over the channel won't panic from a send-on-closed-channel
// race.
func producer(jobStream chan<- *job, terminationSignal chan os.Signal) {

	fmt.Println("starting producer")

	ticker := time.NewTicker(50 * time.Millisecond)
	var counter int

	for {
		select {
		case <-ticker.C:
			counter++
			jobStream <- &job{counter}

		case <-terminationSignal:
			fmt.Println("\nreceived shutdown signal, stopping producer")
			close(jobStream) // close channel on producer side
			return
		}
	}
}
