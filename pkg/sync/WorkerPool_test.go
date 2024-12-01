package sync

import (
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func TestWaitingForFutureToFinish(t *testing.T) {
	pool := NewWorkerPool(2, 0)
	var counter atomic.Int32
	future := NewFutureInPool(pool, func() (string, *error) {
		myTestFunc(&counter)
		return "hello", nil
	})

	future.Wait()

	assert.Equal(t, int32(1), counter.Load())
	assert.Equal(t, "hello", future.Result)
}

func TestWaitingForWorkerpoolToEmpty(t *testing.T) {
	pool := NewWorkerPool(2, 0)
	var counter atomic.Int32
	NewFutureInPool(pool, func() (string, *error) {
		myTestFunc(&counter)
		return "", nil
	})

	pool.Wait()

	assert.Equal(t, int32(1), counter.Load())
}

//
//func TestQueuesFuncWhenPoolIsBusy(t *testing.T) {
//	pool := NewWorkerPool(2, 10)
//
//	var counter atomic.Int32
//	start := time.Now()
//	for _ = range 10 {
//		pool.Exec(func() {
//			myTestFunc(&counter)
//		})
//	}
//
//	pool.Wait()
//	finish := time.Now()
//	log.Print("Elapsed", 250*time.Millisecond > finish.Sub(start))
//	// There are two go routines, each take at least 50ms to run and 10 functions to execute in total so it will take
//	// at least 250ms to run
//	assert.Greater(t, finish.Sub(start), 250*time.Millisecond)
//}
//
//func TestBlocksWhenBufferIsFull(t *testing.T) {
//	pool := NewWorkerPool(2, 1)
//
//	// We have 2 workers, 10 commands and  queue size of 1 so NewFutureInPool will block after the first 2 calls to it
//	var counter atomic.Int32
//	start := time.Now()
//	for _ = range 10 {
//		pool.Exec(func() {
//			myTestFunc(&counter)
//		})
//	}
//	finish := time.Now()
//
//	pool.Wait()
//	assert.Greater(t, finish.Sub(start), 200*time.Millisecond)
//}
//
//func TestClosesGoRoutines(t *testing.T) {
//	before := runtime.NumGoroutine()
//	pool := NewWorkerPool(5, 1)
//	assert.Equal(t, before+5, runtime.NumGoroutine())
//
//	pool.Close()
//	time.Sleep(50 * time.Millisecond)
//	assert.Equal(t, before, runtime.NumGoroutine())
//
//}
//
//func TestPanicsWhenCallExecOnClosedPool(t *testing.T) {
//	pool := NewWorkerPool(5, 1)
//	pool.Close()
//
//	assert.Panics(t, func() {
//		pool.Exec(func() {
//			time.Sleep(50 * time.Millisecond)
//		})
//	})
//}

func myTestFunc(counter *atomic.Int32) {
	time.Sleep(50 * time.Millisecond)
	counter.Add(1)
}
