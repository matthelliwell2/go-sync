package sync

import (
	"sync"
)

// WorkerPool is a pool for Go Routines that can be used for executing Futures.
type WorkerPool struct {
	// The waitgroup used for overall synchronisation of the workerpool, eg waiting until all go routines in the
	// pool have finished.
	wg *sync.WaitGroup
	// The channel used to send functions to the workerpool. As we want to send any function to be executed, we define
	// the channel to take func(). To handle return types etc, the call needs to wrap the call.
	execChan chan func()
	// Whether the workerpool is running
	isRunning bool
}

// NewWorkerPool creates a new workpool
// poolSize - the number of Go Routines in the pool
// bufferSize - The size of the queue used to queue requests to the pool
func NewWorkerPool(poolSize int, bufferSize int) WorkerPool {
	wp := WorkerPool{
		wg:       new(sync.WaitGroup),
		execChan: make(chan func(), bufferSize),
	}

	for range poolSize {
		wp.startWorker()
	}

	wp.isRunning = true

	return wp
}

// NewFutureInPool returns a future that instead of starting the go routine immediately, queues it to a worker pool
// for execution when a worker is free.
//
// wp the workerpool to execute the function
// fn the function to execute in a go routine
//
// Returns a pointer to a Future which can be used to retrieve the functions results when it has executed. If the
// workpool is busy and the buffer is full, this method blocks until a slot becomes free.
func NewFutureInPool[R any, E error](wp WorkerPool, fn func() (R, *E)) *Future[R, E] {
	// Create a future with its own wait group so we can wait on the individual future if we want to. Don't use the
	// NewFuture function as that creates its own go routine and we want to use the worker pool.
	f := Future[R, E]{wg: new(sync.WaitGroup), fn: fn}

	// Increment the wait groups for the future and the work pool so we can wait on either.
	wp.wg.Add(1)
	f.wg.Add(1)

	// The workerpool accepts func() so we wrap the function we actually want to call to capture the results and set them
	// on he future.
	wp.execChan <- func() {
		defer f.wg.Done()

		result, err := f.fn()
		f.Result = result
		f.Error = err
	}

	return &f
}

// Wait waits for all functions either executing or queued to complete. This is equivalent to calling [AllSettled] on
// each of the futures in the pool, but is a bit more convenient.
func (wp WorkerPool) Wait() {
	wp.wg.Wait()
}

// Close closes the worker pool by closing the channel
func (wp WorkerPool) Close() {
	wp.isRunning = false
	close(wp.execChan)
}

// Starts a go routine and then reads from the worker pool channel to execute the required function.
func (wp WorkerPool) startWorker() {
	go func() {
		for fn := range wp.execChan {
			fn()
			wp.wg.Done()
		}
	}()
}
