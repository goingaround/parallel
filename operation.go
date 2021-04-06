package parallel

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// operation is an object built to run multiple anonymous functions simultaneously.
// When invoking Start every function is assigned to a goroutine.
// All goroutines are waiting for the same condition to be set free to execute.
type operation struct {
	mu           sync.Mutex
	funcs        []func()
	once         sync.Once
	ready, ended sync.WaitGroup
	exec, done   chan struct{}
}

func NewOperation(funcs ...func()) (*operation, error) {
	if len(funcs) == 0 {
		return nil, errors.New("got no functions")
	}

	for i, f := range funcs {
		if f == nil {
			return nil, fmt.Errorf("got nil function (index %d)", i)
		}
	}

	return &operation{funcs: funcs}, nil
}

func (op *operation) Run(timeout, delay time.Duration) error {
	op.mu.Lock()
	defer op.mu.Unlock()

	// setup for every run
	op.ready, op.ended, op.exec, op.done =
		sync.WaitGroup{}, sync.WaitGroup{}, make(chan struct{}), make(chan struct{})

	go op.run(delay)

	select { // listen
	case <-op.done:
		return nil
	case <-timeAfter(timeout):
		return &ErrTimeoutExceeded{}
	}
}

func (op *operation) run(delay time.Duration) {
	defer close(op.done)

	// scheduling goroutines
	for pid, f := range op.funcs {
		op.ready.Add(1)
		op.ended.Add(1)
		go op.process(pid, f, delay)
	}

	// waiting for all goroutines to be stuck on wait
	op.ready.Wait()
	// set goroutines free
	close(op.exec)
	// waiting for all goroutines to finish execution
	op.ended.Wait()
}

func (op *operation) process(pid int, f func(), delay time.Duration) {
	op.ready.Done()
	defer op.ended.Done()

	// wait for broadcast signal
	<-op.exec
	// enforce delay
	if delay > 0 {
		time.Sleep(time.Duration(pid) * delay)
	}
	// execution
	f()
}

func timeAfter(timeout time.Duration) <-chan time.Time {
	if timeout <= 0 {
		return make(chan time.Time) // pseudo chan - no timeout
	}
	return time.After(timeout)
}
