// Package parallel provides context-aware helpers for concurrent batch and stream work.
//
// It supports unbounded and bounded task execution, concurrent iteration, and
// order-preserving concurrent mapping and filtering. Fail-fast helpers cancel
// sibling work through context.Context, while MapResults supports best-effort
// processing with one outcome per input value. Pool provides fixed long-lived
// workers and bounded queueing for intermittent background work. [Pool.TrySubmit]
// supports immediate rejection when busy, and [SubmitValue] provides typed futures.
// [MapStream] processes channel inputs with bounded concurrency, per-item results,
// and cancellation shared with producers through [Stream.Context].
package parallel
