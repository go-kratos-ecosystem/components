// Package debounce provides context-aware, trailing-edge debouncing.
//
// [Debouncer.Trigger] retains the latest value until a quiet interval elapses.
// [Debouncer.Flush] forces pending work and [Debouncer.Close] cancels work that
// has not started. A single worker serializes handlers, and [WithOnError] reports
// asynchronous handler failures.
package debounce
