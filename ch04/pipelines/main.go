// Package main demonstrates the Go pipeline pattern: a chain of stages, each
// implemented as a function that takes an input channel and returns an
// output channel, connected by simply passing one stage's output as the
// next stage's input. Every stage runs in its own goroutine, so values flow
// through the pipeline concurrently rather than being fully processed by
// one stage before the next begins. A shared "done" channel gives every
// stage a way to be cancelled early, preventing goroutine leaks if the
// pipeline is torn down before it's fully drained.
package main

import "fmt"

func main() {

	// done is the pipeline's cancellation signal. Closing it (rather than
	// sending on it) lets every stage's goroutine observe the cancellation
	// simultaneously via the same closed-channel receive.
	done := make(chan interface{})
	defer close(done)

	jobs := make([]*job, 0)
	for i := 0; i < 10; i++ {
		jobs = append(jobs, &job{id: i})
	}

	// Wire up the pipeline: generator -> multiply(x2) -> add(+2) -> multiply(x50).
	// Each stage is a separate goroutine reading from the previous stage's
	// output channel and writing to its own.
	jobStream := generator(done, jobs...)
	pipeline := multiply(done, add(done, multiply(done, jobStream, 2), 39), 10)

	// Draining the final stage's channel is what actually pulls values
	// through every upstream stage.
	for result := range pipeline {
		fmt.Println(result)
	}
}

// job is the value type flowing through every stage of the pipeline.
type job struct {
	id int
}

// generator is the pipeline's source stage: it feeds the pre-built jobs
// slice into a channel, one at a time, in its own goroutine so callers can
// start consuming before all jobs have been sent.
func generator(done chan interface{}, jobs ...*job) <-chan *job {
	jobStream := make(chan *job)

	go func() {
		// Closing jobStream when this goroutine exits (whether the loop
		// finished normally or was cancelled) signals downstream stages
		// that no more values are coming, letting their range loops end.
		defer close(jobStream)
		for _, job := range jobs {

			// select on done vs. the send lets a blocked send be abandoned
			// immediately if the pipeline is cancelled, instead of leaking
			// this goroutine forever waiting for a receiver that will
			// never come.
			select {
			case <-done:
				return
			case jobStream <- job:
			}
		}
	}()

	return jobStream
}

// multiply is a transform stage that multiplies every job's id by value.
// The same done/select cancellation pattern as generator is used both for
// the outgoing send and, implicitly, by ranging over jobStream (which ends
// once the upstream stage closes its channel).
func multiply(done chan interface{}, jobStream <-chan *job, value int) <-chan *job {
	stream := make(chan *job)

	go func() {
		defer close(stream)
		for j := range jobStream {
			select {
			case <-done:
				return
			case stream <- &job{id: j.id * value}:
			}
		}
	}()

	return stream
}

// add is a transform stage that adds value to every job's id. Structurally
// identical to multiply - each pipeline stage follows the same shape:
// range over input, transform, select-send to output, honoring done.
func add(done chan interface{}, jobStream <-chan *job, value int) <-chan *job {
	stream := make(chan *job)

	go func() {
		defer close(stream)
		for j := range jobStream {
			select {
			case <-done:
				return
			case stream <- &job{id: j.id + value}:
			}
		}
	}()

	return stream
}
