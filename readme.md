# Async

A collection of async utilities for Go.

The sync package provides methods for synchronising GoRoutines, such as waiting until all GoRoutines complete. You have the option of spawning as many GoRoutines as you need to invoke or you can limit their number using a WorkerPool. A summary of the methods is in the package documentation. The documentation for each method contains examples on how to use it.

## Future
Like a future or promise in other languages, a future runs a task asynchronously and lets you wait until the task is complete. NewFuture creates a new GoRoutines and executes the provided function in that GoRoutine. Use the wait method to wait for the GoRoutine to complete.

Futures can be combined together is various combinations to make them more useful, such as waiting until they all finish.

Example:

```go
f := NewFuture[string, error](func() (string, *error) {
time.Sleep(time.Second)
return "Hello world", nil
})

fmt.Println("Start", time.Now())
result, err := f.Wait()
fmt.Println("End  ", time.Now())
fmt.Println("Results", result, err)
```

Output:

    Start 2024-06-15 17:00:24.9777397 +0100 BST m=+0.003647301
    End   2024-06-15 17:00:25.9870966 +0100 BST m=+1.013004201
    Results Hello world <nil>

## AllSettled

This function takes multiple futures and returns a future that completes when all the input futures have succeeded or failed.

For example, if you have a list of insert operations that you want to run on a database and you want them to run in parallel, you could use AllSettled to run them in parallel and check the results of each operation:

```Go
data := []string{"value1", "value2", "value3"}
futures := make([]*Future[int, error], len(data))
for i, value := range data {
    f := NewFuture(func() (int, *error) {
        fmt.Println("Calling db with value", value)
        // Do some actual database calls here
        return i, nil
    })
    futures = append(futures, f)
}

results, errors := AllSettled(futures...).Wait()

// Check the results and errors
```

## Any

This function takes multiple futures and returns a future that completes when any the input futures have succeeded or they have all failed.

## Worker Pool

When you create a future using NewFuture, it runs in its own Go Routine. In some circumstances, you might want to limit the number of concurrent futures. For example, if you are sending a large number of inserts to a database, you might need to limit the number running at the same time to avoid overloading the database. For this you can use a worker pool.

A worker pool is a fixed number of Go Routines and a fixed sized queue. When you use a worker pool to execute a future, it first checks if there are any available Go Routines in the pool. If there are any free, the future starts executing immediately. If there aren't any free, then the future is added to a queue. It will be pulled from the queue when a Go Routine becomes available.

If the queue becomes full, then the call to add a future to the worker pool will block until some futures are consumed from the queue.

To create a worker pool, use
```go
NewWorkerPool
```
To execute a future in the worker pool, instead of in a new Go Routine, use
```go
NewFutureInPool
```

Calls to AllSettled etc aren't affected by whether the input futures are in a worker pool or not.

## TODO

* Add examples for everything
* All - TODO
* race - TODO
* each - TODO
* parallelmap - Add outline readme
* groupby - TODO
* should any of this be using iterators?
