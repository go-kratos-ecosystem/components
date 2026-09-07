# Batcher

`batcher` combines individual items into batches for a serial processing function.
Batches are submitted when full, when their accumulation interval expires, or when
explicitly flushed. A bounded queue applies backpressure to producers.

## Installation

```bash
go get github.com/go-fries/fries/batcher/v4
```

## Process items in batches

Provide a processing context whose lifetime covers the batcher's work:

```go
processingContext, cancelProcessing := context.WithCancel(context.Background())
defer cancelProcessing()

b, err := batcher.New(processingContext,
	func(ctx context.Context, items []Record) error {
		return repository.InsertBatch(ctx, items)
	},
	batcher.WithBatchSize(100),
	batcher.WithFlushInterval(time.Second),
	batcher.WithQueueSize(1000),
)
if err != nil {
	return err
}
```

Submit items from one or more producers:

```go
if err := b.Add(ctx, record); err != nil {
	return err
}
```

`Add` returns once the item has been accepted. If the queue is full, it waits
until capacity becomes available, its context is canceled, or shutdown begins.
Canceling an `Add` context after acceptance does not cancel the item.

The worker processes batches serially in queue acceptance order. The interval
starts when the worker receives the first item of an empty batch; adding more
items does not restart it. Queue waiting and handler execution can add latency.
Empty batches are never sent to the handler.

Each handler call receives a separate slice that the batcher does not reuse.
Items are not deep-copied: producers must synchronize mutations to maps, slices,
pointers, or other referenced data after submission.

## Reject immediately when busy

Use `TryAdd` for optional work that should not wait for queue capacity:

```go
err := b.TryAdd(record)
if errors.Is(err, batcher.ErrFull) {
	return nil // Skip this optional item.
}
if err != nil {
	return err
}
```

Rejected items never reach the handler. With an unbuffered queue, `TryAdd`
succeeds only when the worker is ready to receive.

## Flush accepted work

```go
if err := b.Flush(ctx); err != nil {
	return err
}
```

`Flush` places a barrier in the same FIFO queue as items. It waits for all items
before that barrier, including all items accepted before the call, and for their
result callbacks. A partial batch is processed immediately. Later items do not
delay completion of this barrier.

Canceling the flush context stops admission or waiting. An already accepted
barrier continues in the background. After shutdown begins, `Flush` returns
`ErrClosed`; use `Shutdown` to wait for the final outcome.

## Observe processing results

Register a callback to observe each batch, including successful ones:

```go
batcher.WithOnResult(func(result batcher.Result) {
	if result.Err != nil {
		log.Printf("batch of %d items failed: %v", result.Size, result.Err)
	}
})
```

Each batch is passed to the handler once. A handler error does not stop later
batches or trigger an automatic retry. Implement any retry or failed-item
handling inside the handler, where the items are available. The handler's error
alone does not tell the batcher whether the destination partially accepted data.

Errors are cumulative: `Flush` reports the first handler error since construction
through its barrier, and `Shutdown` reports the first handler error over the
entire lifetime. Reading an error does not clear it. This keeps error storage
bounded and prevents a canceled or concurrent flush from consuming another
caller's error. Use `WithOnResult` to observe subsequent failures.

The handler and result callback run on the same worker. They must not
synchronously call blocking `Add`, `Flush`, or `Shutdown` methods on that batcher.
Callbacks contribute to backpressure. Panics follow normal Go semantics and are
not recovered.

## Shut down

```go
shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

if err := b.Shutdown(shutdownContext); err != nil {
	return err
}
```

`Shutdown` stops admission, unblocks waiting producers, and drains accepted items,
including the last partial batch. Submissions already in progress may succeed
concurrently with shutdown; they are included in the drain. Repeated shutdown
calls wait for the same final result.

Canceling the shutdown context only stops that wait. To cancel processing, cancel
the context supplied to `New`. That also starts shutdown: remaining accepted
items still reach the handler, but with the canceled processing context. Handlers
should return promptly when it is canceled. Handler return values determine the
processing outcome, even if cancellation races with successful completion.

Items are held in memory until processing and are not durable across process
exit. `Add` succeeding confirms acceptance into memory, not persistence at the
destination.

## Options

| Option | Default | Meaning |
| --- | --- | --- |
| `WithBatchSize` | `100` | Maximum items in a batch; must be positive. |
| `WithFlushInterval` | `1s` | Maximum accumulation interval once a batch starts; must be positive. |
| `WithQueueSize` | `1000` | Pending items or flush barriers; zero is unbuffered. |
| `WithOnResult` | `nil` | Synchronous callback after each handler call. |

Queue capacity excludes the batch currently accumulating or executing. The
batcher holds at most queue capacity plus batch size items internally, apart
from values retained by callers. Producers blocked in `Add` also retain their
own items, so callers should bound producer concurrency when necessary.

Options are applied in order, with the last value taking effect. `New` rejects
invalid final sizes or intervals with an error wrapping `ErrInvalidOption` and
ignores nil options.
