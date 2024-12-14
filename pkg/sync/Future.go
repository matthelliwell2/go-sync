// Package sync A collection of async utilities for Go.
// The package provides methods for synchronising GoRoutines, such as waiting until all GoRoutines complete. You have
// the option of spawning as many GoRoutines as you need to invoke or you can limit their number using a WorkerPool. A
// summary of each method is below. Refer to the docs for each method for details and examples.
//
// #[Future]
//
// The basic type used by this package. It represents a function that will complete at some point in the future. Use
// [NewFuture] or [NewFutureInPool] to create one and then [All], [AllSettled] etc to combine them.
//
// #[All]
//
// Waits for all the specified futures to complete successfully or for one of them to fail.
//
// #[AllSettled]
//
// This function takes multiple futures and returns a future that completes when all the input futures have succeeded or
// failed.
//
// #[Any]
//
// This function takes multiple futures and returns a future that completes when any the input futures have succeeded or
// they have all failed.
//
// #[ParallelMap]
//
// Executes a function in parallel for each element of a slice.
//
// #[WorkerPool]
//
// Allows futures to run a pool of GoRoutines, instead of creating a new GoRoutine for each future.
package sync

import (
	"sync"
)

// Future represents a function that will complete and return a value or error at some point in the future.
// T is return type of the function
// E is the error type of the function
type Future[T any, E error] struct {
	// The result of executing fn in a Go Routine
	Result T
	// The error result from the function. We are using a pointer because we don't know what type E is so the compiler
	// won't let us check for nil.
	Error *E
	// The wait group used internally to wait on the function completing
	wg *sync.WaitGroup
	// The function to be called. If you want to call a function that takes parameters, you will need to wrap it in a
	// closure. Note that the function must return *E and not E. This is so that the package can check if a nil value
	// is returned.
	fn func() (T, *E)
}

// NewFuture creates a new future for the specified function and executes it in a Go Routine.
// T is return type of the function.
// E is the error type of the function.
// fn The function to be executed in the Go Routine. Note that it returns (T, *E), not (T, E).
// Returns a pointer to a [Future]. The caller can wait directly on this future or combine the future with other futures
// using [AllSettled] etc. A pointer is returned so that the caller and the Go Routine both reference the same
// object, for example so the Go Routine can update the value of the result and the caller can read that value.
//
// For example
//
//		f := NewFuture[string, error](func() (string, *error) {
//			time.Sleep(time.Second)
//			return "Hello world", nil
//		})
//
//		fmt.Println("Start", time.Now())
//		result, err := f.Wait()
//		fmt.Println("End  ", time.Now())
//	 fmt.Println("Results", result, err)
//
// Output:
//
//	Start 2024-06-15 17:00:24.9777397 +0100 BST m=+0.003647301
//	End   2024-06-15 17:00:25.9870966 +0100 BST m=+1.013004201
//	Results Hello world <nil>
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
