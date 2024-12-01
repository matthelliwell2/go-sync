# Async

A collection of async utilities for Go.

## Future
Like a future or promise in other languages, a future runs a task asynchronously and lets you wait until the task is complete.

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

Futures can be combined together is various combinations to make them more useful

## All

This takes multiple futures and returns a future that completes when all the input futures have succeeded or any single future fails.

## AllSettled

This takes multiple futures and returns a future that completes when all the input futures have succeeded or failed.

## Any

This takes multiple futures and returns a future that completes when any the input futures have succeeded or they have all failed.

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

## Worker Pool

When you create a future using NewFuture, it runs in its own Go Routine. In some circumstances, you might want to limit the number of concurrent futures. For example, if you are sending a large number of inserts to a database, you might need to limit the number running at the same time to avoid overloading the database. For this you can use a worker pool.

A worker pool is a fixed number of Go Routines and a fixed sized queue. When you use a worker pool to execute a future, it first checks if there are any available Go Routines in the pool. If there are any free, the future starts executing immediately. If there aren't any free, then the future is added to a queue. It will be pulled from the queue when a Go Routine becomes available.

If the queue becomes full, then the call to execute a future in the worker pool will block until some futures are consumed from the queue.

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

* Future - done
* Workerpool - done
* AllSettle - done
* All - TODO
* race - TODO
* each - TODO
* map - TODO
* groupby - TODO
* should any of this be using itertors?
