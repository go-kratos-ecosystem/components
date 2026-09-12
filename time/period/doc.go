// Package period provides immutable, half-open time intervals and set operations.
//
// A [Period] includes its start and excludes its end: [start, end). Equal
// endpoints describe an empty interval. [New] normalizes endpoints to UTC and
// removes monotonic clock readings, so all operations compare absolute times.
//
// Use [Period.Intersect] and [Period.Subtract] for individual intervals, [Merge]
// to normalize coverage, and [Gaps] to find missing coverage inside its extent.
// Operations are synchronous and do not modify their inputs.
package period
