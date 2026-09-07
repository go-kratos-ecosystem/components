# Parallel

`parallel` provides small, context-aware helpers for concurrent batch and stream work.
It propagates errors and cancellation, supports explicit concurrency limits,
and provides an optional fixed-worker pool for intermittent background work.

## Installation

```bash
go get github.com/go-fries/fries/parallel/v4
```

## Run tasks

```go
err := parallel.RunLimit(ctx, 4,
	func(ctx context.Context) error {
		return refreshCache(ctx)
	},
	func(ctx context.Context) error {
		return updateIndex(ctx)
	},
)
```

The first task error cancels the context passed to the remaining tasks. `Run`
provides the same behavior without a concurrency limit.

## Process collections

Use `ForEach` for concurrent side effects:

```go
err := parallel.ForEach(ctx, 8, userIDs, func(ctx context.Context, id int64) error {
	return notifyUser(ctx, id)
})
```

Use `Map` for type-safe transformations. Results preserve input order even when
callbacks finish in a different order:

```go
profiles, err := parallel.Map(ctx, 8, userIDs,
	func(ctx context.Context, id int64) (Profile, error) {
		return loadProfile(ctx, id)
	},
)
```

## Keep partial results

Use `MapResults` when every item should be attempted even if some callbacks
fail. Each result corresponds to the input at the same index:

```go
results, batchErr := parallel.MapResults(ctx, 8, userIDs,
	func(ctx context.Context, id int64) (Profile, error) {
		return loadProfile(ctx, id)
	},
)

for index, result := range results {
	if result.Err != nil {
		log.Printf("load user %d: %v", userIDs[index], result.Err)
		continue
	}
	useProfile(result.Value)
}

if batchErr != nil {
	return batchErr
}
```

## Filter values

`Filter` evaluates a predicate concurrently and preserves input order:

```go
activeUsers, err := parallel.Filter(ctx, 8, users,
	func(ctx context.Context, user User) (bool, error) {
		return service.IsActive(ctx, user.ID)
	},
)
```

## Process streaming inputs

Use `MapStream` when inputs arrive over time or results should be consumed as
they become available. It starts a fixed number of workers and returns a stream
of indexed, per-item results:

```go
input := make(chan int64)
stream, err := parallel.MapStream(ctx, 8, input,
	func(ctx context.Context, id int64) (Profile, error) {
		return loadProfile(ctx, id)
	},
)
if err != nil {
	return err
}
defer stream.Close()

go func() {
	defer close(input)
	for _, id := range userIDs {
		select {
		case input <- id:
		case <-stream.Context().Done():
			return
		}
	}
}()

for result := range stream.Results() {
	if result.Err != nil {
		log.Printf("input %d failed: %v", result.Index, result.Err)
		continue
	}
	useProfile(result.Value)
}

if err := stream.Wait(ctx); err != nil {
	return err
}
```

Results are delivered as workers make them available, without input-order
sorting. `Index` is the zero-based input receive position. A slow earlier input
does not hold back a ready result from another worker. Concurrently ready
results have no deterministic ordering.

Each worker receives another input only after handing off its previous result.
The result channel is unbuffered and there is no internal prefetch queue, so a
slow consumer slows input consumption. At most the concurrency limit's worth of
items are being received, processed, or delivered inside the stream. Producer
buffers and data retained by callbacks or consumers are additional.

Callback errors are included in individual results, alongside any returned
value, and other inputs continue processing. `Wait` returns nil after normal
completion even if some callbacks failed. Parent cancellation or an early
`Close` is reported as the stream's cancellation cause. An already canceled
parent context, a non-positive limit, a nil input, or a nil callback causes
construction to fail before any worker starts.

### Stop consumption early

When enough results have arrived, call `Close`, or return from a function with
`defer stream.Close()` registered. Breaking out of the result loop alone does
not cancel the stream.

`Close` cancels the shared stream context and waits for internal workers. It
unblocks input reception and result delivery, but running callbacks still need
to observe their context and return. Callbacks must not synchronously call
`Close` or `Wait` on their own stream. Panics follow normal Go semantics and are
not recovered.

Producers should select on `stream.Context().Done()` while sending, as shown
above. The caller owns the input channel and any producer goroutines: the stream
neither closes that channel nor waits for those goroutines. Join the producer
separately when its resource cleanup must finish before returning.

On cancellation, received inputs may have no delivered result; the stream does
not synthesize results for unread inputs. Already delivered results remain
usable. Inputs and results are not deep-copied, so referenced data requires
appropriate synchronization if mutated.

`Wait` observes completion and does not drain results. Consume results or cancel
the stream before waiting indefinitely. Its context controls only that wait;
canceling it does not cancel processing. The shared stream context is also
canceled on normal completion, so use `Wait` with an independent context to
obtain the terminal status. Repeated waits return the same terminal result.

## Reuse fixed workers

Use `Pool` for intermittent work that should share a fixed concurrency limit
and a bounded queue:

```go
pool := parallel.NewPool(8, parallel.WithQueueSize(32))

future, err := pool.Submit(taskContext, func(ctx context.Context) error {
	return refreshCache(ctx, cacheKey)
})
if err != nil {
	return err
}

// Wait only when this request needs the task result.
if err := future.Wait(ctx); err != nil {
	return err
}
```

`Submit` returns after the task is accepted; execution may begin immediately or
after an earlier task finishes. `Execute` combines submission and waiting for
synchronous handlers.

The context passed to `Submit` is also passed to the task. If background work
must outlive a request, create an explicit task context, such as one derived
with `context.WithoutCancel`, instead of using the request context directly.
Task panics follow normal Go semantics and are not recovered by the pool.

### Skip work when the pool is busy

`TrySubmit` returns immediately with `ErrPoolFull` when the task cannot be
accepted. Rejected tasks never run. With an unbuffered queue, submission succeeds
only when a worker is ready to receive the task.

```go
_, err := pool.TrySubmit(taskContext, func(ctx context.Context) error {
	return refreshCache(ctx, cacheKey)
})
if errors.Is(err, parallel.ErrPoolFull) {
	return nil // Skip this optional refresh while the pool is busy.
}
if err != nil {
	return err
}
```

### Retrieve a typed result

`SubmitValue` uses the same queue and backpressure as `Submit`, returning a
`ValueFuture[T]` that carries the callback's value and error:

```go
future, err := parallel.SubmitValue(taskContext, pool,
	func(ctx context.Context) (Profile, error) {
		return loadProfile(ctx, userID)
	},
)
if err != nil {
	return err
}

// Other work can run before waiting for the result.
profile, err := future.Wait(ctx)
if err != nil {
	return err
}
useProfile(profile)
```

Multiple callers can wait on the same future. A completed future preserves both
the value and error, including values returned alongside an error. Canceling a
wait returns the zero value and cancellation cause without canceling the task.
If cancellation prevents the task from starting, its result is the zero value
and the task context's cancellation cause. Referenced result data is shared;
callers must synchronize mutations.

### Shut down the pool

Shut the pool down when its owner stops:

```go
shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

if err := pool.Shutdown(shutdownContext); err != nil {
	return err
}
```

If the shutdown context expires, `Shutdown` returns its cancellation cause but
the pool continues draining accepted tasks in the background. Task contexts
remain responsible for canceling work that should not outlive shutdown.

Batch helpers wait for started work to return. Pool submission returns after
acceptance and exposes completion through `Future`. Callbacks and tasks should
observe the provided context and stop promptly after cancellation.
