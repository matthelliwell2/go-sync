package sync

// Any returns a future which completes when any of the input futures complete successfully.
//
// It will complete as soon as one of the futures completes, it will not wait on all the others. If they all fail then
// the error is an array of errors for each function. In this case the result is undefined and shouldn't be used.
func Any[T any, E error](futures ...*Future[T, E]) *Future[T, Errors[E]] {
	return NewFuture(func() (T, *Errors[E]) {
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

		// Wait on results for the channel we've either got all error or we've got the first successful result
		errors := make([]*E, len(futures))
		for range len(futures) {
			answer := <-q
			if answer.err == nil {
				return answer.result, nil
			}
			// Collect the errors in case all the futures fail and we need to return them. The errors are collected in the
			// same order as the function args, not in the order in which they completed.
			errors[answer.index] = answer.err
		}

		err := Errors[E]{Errors: errors}
		var dummyResult T
		return dummyResult, &err
	})
}
