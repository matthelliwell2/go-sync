package sync

// All returns a future that completes either when all the futures complete or any future fails. If any of the input
// future fails, the returned future will contain nil for the results and the error of the first failed future.
// Otherwise the future has an array of each result. The future will complete when the first failure is detected, rather
// than waiting for all future to complete before checking if any have failed.
// If there is a failure, the functions will continue to execute in their go routines.
func All[T any, E error](futures ...*Future[T, E]) *Future[[]T, E] {
	return NewFuture(func() ([]T, *E) {
		// Make a channel to collect the results of the futures. Make the channel big enough to hold all the results
		// so even if this routine completes and nothing is reading from the channel, the futures will still run and
		// complete without blocking. This allows the memory for the channel etc to be garbage collected.
		q := make(chan resultWithIndex[T, E], len(futures))

		// Get each future to put its results on a channel without waiting for any others to finish so we can get the
		// result of the first one.
		for i, future := range futures {
			go func() {
				result, err := future.Wait()
				q <- resultWithIndex[T, E]{result: result, err: err, index: i}
			}()
		}

		results := make([]T, len(futures))
		for range len(futures) {
			answer := <-q
			if answer.err != nil {
				return nil, answer.err
			}
			results[answer.index] = answer.result
		}

		return results, nil
	})
}
