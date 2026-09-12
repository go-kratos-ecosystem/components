# Period

`period` provides immutable time intervals and operations for finding overlaps,
remaining availability, combined coverage, and gaps.

Every interval is half-open: **`[start, end)`** includes its start and excludes
its end. Two reservations from 09:00–10:00 and 10:00–11:00 are adjacent and do
not overlap. Merging them produces 09:00–11:00.

## Installation

```bash
go get github.com/go-fries/fries/period/v4
```

The module uses only the Go standard library at runtime.

## Find available time

Subtract reservations from an available window:

```go
package main

import (
	"fmt"
	"time"

	"github.com/go-fries/fries/period/v4"
)

func main() {
	start := time.Date(2026, time.September, 12, 9, 0, 0, 0, time.UTC)
	available, err := period.New(start, start.Add(9*time.Hour))
	if err != nil {
		panic(err)
	}
	morning, err := period.New(start.Add(time.Hour), start.Add(2*time.Hour))
	if err != nil {
		panic(err)
	}
	afternoon, err := period.New(start.Add(5*time.Hour), start.Add(7*time.Hour))
	if err != nil {
		panic(err)
	}

	for _, free := range available.Subtract(morning, afternoon) {
		fmt.Printf("%s–%s\n", free.Start().Format("15:04"), free.End().Format("15:04"))
	}
}
```

Output:

```text
09:00–10:00
11:00–14:00
16:00–18:00
```

`Subtract` accepts unsorted, duplicate, overlapping, adjacent, and out-of-window
cuts. It returns sorted, non-empty intervals containing all remaining time.

## Combine coverage and find gaps

For data received in several time windows, merge the windows to normalize the
coverage, or find the missing intervals between them:

```go
covered := period.Merge(received...)
missing := period.Gaps(received...)
```

```text
Received:  [00:00, 02:00), [01:00, 03:00), [05:00, 06:00)
Covered:   [00:00, 03:00), [05:00, 06:00)
Missing:   [03:00, 05:00)
```

`Gaps` only reports gaps between the earliest start and latest end of non-empty
input intervals. To include missing time at the edges of a requested window,
use `requested.Subtract(received...)`.

## API

| Operation | Behavior |
| --- | --- |
| `New(start, end)` | Constructs an interval; returns `ErrInvalidRange` when end precedes start. |
| `Start()` / `End()` | Returns the endpoints in UTC. |
| `IsEmpty()` | Reports whether the endpoints are equal. |
| `Duration()` | Returns elapsed time, with the saturation behavior of `time.Time.Sub`. |
| `Equal(other)` | Compares the represented sets of time points. All empty intervals are equal. |
| `Contains(t)` | Includes the start, excludes the end. |
| `ContainsPeriod(other)` | Reports whether every point in other is contained. Every interval contains every empty interval. |
| `Overlaps(other)` | Reports whether there is a non-empty intersection. |
| `Adjacent(other)` | Reports whether two non-empty intervals meet at an endpoint. |
| `Intersect(other)` | Returns the shared interval and true, or the zero value and false. |
| `Subtract(others...)` | Removes cuts and returns sorted remaining intervals. |
| `Merge(periods...)` | Combines overlapping and adjacent intervals into sorted coverage. |
| `Gaps(periods...)` | Returns sorted internal gaps in the combined coverage. |

## Empty intervals

Equal endpoints are valid: `New(t, t)` produces an empty interval at `t`.
The zero value of `Period` is also empty and has zero-time endpoints.

- Empty intervals contain no time points, overlap nothing, and are never adjacent.
- Every interval contains every empty interval, regardless of its endpoints.
- `Equal` treats all empty intervals as the same empty set. Use `Equal` when
  comparing interval meaning; Go's `==` compares the stored endpoints.
- Set operations ignore empty inputs and return nil slices for empty results.

## Time and ownership semantics

`New` converts endpoints to UTC and strips monotonic clock readings before
validation. All comparisons use absolute time with nanosecond precision.
To display results in another location, use `p.Start().In(location)` and
`p.End().In(location)`.

Construct calendar boundaries in the desired location before calling `New`.
For example, local midnight to the next local midnight can span 23 or 25 hours
across daylight-saving changes. `Duration` reports that actual elapsed time.
For intervals longer than `time.Duration` can represent, it saturates at the
maximum duration; interval comparisons and set operations still use endpoints.

Values are immutable and safe to copy or read concurrently. Operations never
modify or reuse caller-owned input slices. `Merge`, `Gaps`, and `Subtract` use
O(n log n) time and O(n) additional space for n input intervals or cuts. Results
are sorted, non-empty, and neither overlapping nor adjacent.
