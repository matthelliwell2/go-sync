package sync

// ParallelMap iterates over all elements of a slice, calling a function on each one and capturing the result in another
// slice. Each function call is executed in a Go Routine so they run in parallel.
//
// S - the type of the elements of the input slice
// R - the type of the elements of the resulting slice
//
// Returns a future containing the results of the calling the function on each element. If any function calls return an error,
// the future completes early and no results are returned.
// TODO add option to run these in a workerpool
func ParallelMap[S any, R any, E error](slice []S, fn func(S) (R, *E)) *Future[[]R, E] {
	var futures = make([]*Future[R, E], len(slice))
	for i, e := range slice {
		futures[i] = NewFuture(func() (R, *E) {
			return fn(e)
		})
	}

	return All(futures...)
}
