# Debounce

`debounce` retains the latest value and invokes a handler after triggers have
been quiet for a configured interval. One worker serializes handler calls.

For a 500 ms delay:

```text
  0 ms  Trigger("a")    -> due at 500 ms
200 ms  Trigger("ab")   -> due at 700 ms
400 ms  Trigger("abc")  -> due at 900 ms
900 ms  handler("abc")
```

Continuous triggers keep postponing execution. Each debouncer manages one
latest-value slot; use separate instances for independent actions, such as
saving different drafts.

## Installation

```bash
go get github.com/go-fries/fries/debounce/v4
```

## Save the latest value

Create a debouncer with a processing context that covers its lifetime:

```go
saver, err := debounce.New(ctx, 500*time.Millisecond,
	func(ctx context.Context, content string) error {
		return repository.SaveDraft(ctx, draftID, content)
	},
	debounce.WithOnError(func(err error) {
		log.Printf("save draft: %v", err)
	}),
)
if err != nil {
	return err
}
defer saver.Close()
```

Whenever the content changes, trigger a save:

```go
if err := saver.Trigger(content); err != nil {
	return err
}
```

`Trigger` returns once the value is stored. It replaces the previous pending
value and restarts the delay. Success means acceptance, not that the value has
been saved; another trigger can replace it before execution.

Handler calls never overlap. Triggers received during a running handler form
the next invocation. If that invocation's quiet interval has already elapsed
when the handler finishes, it can start immediately. Otherwise it waits for the
remaining interval. Handler execution and scheduling can add latency.

Values are not deep-copied. Synchronize mutations to maps, slices, pointers, or
other referenced data that a handler may read after submission.

## Save immediately

Use `Flush` when a user explicitly saves or before closing the debouncer:

```go
if err := saver.Flush(waitContext); err != nil {
	return err
}
```

`Flush` freezes the latest pending value and schedules it without the remaining
quiet interval. It waits for that invocation and its error callback. Later
triggers form a new pending value and cannot replace the frozen one.

If there is no pending value, `Flush` waits for an existing frozen or running
invocation, or returns nil if idle. It returns the selected handler's error;
errors from already completed invocations are not replayed. Concurrent flushes
may wait for the same invocation and receive the same result.

The debouncer holds at most one running value, one frozen value, and one
replaceable pending value. If another frozen invocation already occupies the
slot and a newer value is pending, `Flush` first waits for that frozen invocation
to finish, then captures the latest available pending value. During this
admission wait, pending values may still be replaced or processed by the worker.

Canceling the flush context ends only admission or waiting. An already frozen
invocation continues with the processing context provided to `New`. The result
of an invocation already completed while waiting takes precedence over wait
cancellation. Nil contexts are rejected with `ErrInvalidContext`.

## Handle errors

`WithOnError` reports every handler error, including failures from invocations
initiated by `Flush`. It does not run for successful calls or for values that
were discarded without invoking the handler.

Handler failures do not stop later invocations and are not automatically
retried. Register `WithOnError` to observe background failures; without it,
handler errors are available only to flush callers waiting for that invocation.

The error callback runs on the worker, so `Flush` and `Close` wait for it to
finish. Handlers and error callbacks may call `Trigger`, but must not
synchronously call `Flush` or `Close` on their own debouncer. Panics follow
normal Go semantics and are not recovered.

## Close

```go
saver.Close()
```

`Close` rejects new triggers and flushes, discards pending and frozen values,
cancels the processing context, and waits for the running handler and worker.
It is safe to call repeatedly or concurrently. Handlers must observe their
context and return promptly for closure to finish promptly.

To preserve the final pending value, flush while the processing context is still
active, then close:

```go
err := saver.Flush(waitContext)
saver.Close()
return err
```

Canceling the context supplied to `New` also starts closure. A handler already
selected by the worker counts as running and may receive a canceled context.
Flush callers waiting on a frozen invocation discarded by closure receive
`ErrClosed`. Flush callers waiting on a running invocation receive that handler's
result, unless their own wait context ends first.

The debouncer owns no persistent storage. Closing a debouncer is a cancellation
operation; processing the final value is an explicit `Flush` operation.

## Configuration

`New` requires a positive delay, a non-nil handler, and a non-nil, active context.
Invalid arguments return `ErrInvalidDelay`, `ErrNilHandler`, `ErrInvalidContext`,
or the context's cancellation cause before starting a worker.

`WithOnError` is optional. Options are applied in order; the last callback wins.
A nil callback disables notifications, and nil options are ignored.
