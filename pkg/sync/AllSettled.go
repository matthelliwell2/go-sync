package sync

// AllSettled returns a future that wait on all the input futures to complete, whether they are successful or a
// failure. All the futures must be the same type, otherwise the generic types wpn't work. If you need to wait on
// futures of different types, you can call AllSettleed multiple times. Alternatively, if all the futures are in the
// same worker pool, you can wait on the worker pool and then check the result of each future.
// If no functions fail, an array of results and nil. If one or more functions fail, it returns an array of errors,
// some of which may be nil, depending on whether which functions returns errors.
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
