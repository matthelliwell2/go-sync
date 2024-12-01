package sync

import (
	"sync"
)

type Future[T any, E error] struct {
	// The result of executing fn in a Go Routine
	Result T
	// The error result from the function. We are using a pointer because we don't know what type E is so the compiler won't
	// let us check for nil.
	Error *E
	// The wait group used internally to wait on the function completing
	wg *sync.WaitGroup
	// The function to be called. If you want to call a function that takes parameters, you will need to wrap it in a closure.
	fn func() (T, *E)
}

// NewFuture creates a new future for the specified function and executes it in a Go Routine.
// return a pointer to a future otherwise it will create a copy of the future for the return value but the GO Routines
// will try and update the original future.
// The async function returns a pointer to something that implements Error so that we can check for nil.
func NewFuture[T any, E error](fn func() (T, *E)) *Future[T, E] {
	f := Future[T, E]{wg: new(sync.WaitGroup), fn: fn}
	f.wg.Add(1)
	go func() {
		defer f.wg.Done()

		result, err := f.fn()
		f.Result = result
		f.Error = err
	}()

	// Return a pointer to the future and not a copy so the version being updated by the go routine is the same version
	// seen by the caller.
	return &f
}

// Wait waits for the function to finish.
func (f *Future[T, E]) Wait() (T, *E) {
	f.wg.Wait()
	return f.Result, f.Error
}
