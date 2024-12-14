package sync

// AllSettled returns a future that wait on all the input futures to complete, whether they are successful or a
// failure.
//
// All the futures must be the same type, otherwise the generic types will not work. If you need to wait on
// futures of different types, you can call AllSettled multiple times. Alternatively, if all the futures are in the
// same worker pool, you can wait on the worker pool and then check the result of each future.
// If no functions fail, the future will return an array of results and nil for the error. If one or more functions fail,
// it returns nil for the results and an array of errors, some of which may be nil, depending on whether which functions
// returned errors.
//
// The order of the results and error arrays is the same as the order of the functions passed in.
func AllSettled[T any, E error](futures ...*Future[T, E]) *Future[[]T, Errors[E]] {
	return NewFuture(func() ([]T, *Errors[E]) {
		results := make([]T, len(futures))
		errors := make([]*E, len(futures))
		foundErrors := false
		for index, future := range futures {
			results[index], errors[index] = future.Wait()
			foundErrors = foundErrors || errors[index] != nil
		}

		if foundErrors {
			err := Errors[E]{Errors: errors}
			return results, &err
		} else {
			return results, nil
		}
	})
}
