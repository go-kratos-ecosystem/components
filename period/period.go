package period

import "time"

// Period is an immutable half-open interval [start, end).
//
// Its start is included and its end is excluded. Equal endpoints describe an
// empty interval; the zero value is the empty interval at the zero time.
// Construct other intervals with [New]. Values can be copied and read
// concurrently without synchronization.
//
// Endpoints use UTC with no monotonic clock reading. Period operates on absolute
// times without rounding or calendar arithmetic. Use [Period.Equal] for set
// equality: empty intervals at different endpoints represent the same set.
type Period struct {
	start time.Time
	end   time.Time
}

// New returns the interval [start, end), normalizing both endpoints to UTC and
// removing monotonic clock readings before comparing them.
//
// Equal endpoints are valid and describe an empty interval. If end precedes
// start, New returns the zero value and [ErrInvalidRange].
func New(start, end time.Time) (Period, error) {
	start = start.UTC()
	end = end.UTC()
	if end.Before(start) {
		return Period{}, ErrInvalidRange
	}

	return Period{start: start, end: end}, nil
}

// Start returns the included start endpoint in UTC.
func (p Period) Start() time.Time {
	return p.start
}

// End returns the excluded end endpoint in UTC.
func (p Period) End() time.Time {
	return p.end
}

// IsEmpty reports whether p contains no time points.
func (p Period) IsEmpty() bool {
	return p.start.Equal(p.end)
}

// Duration returns the elapsed time from start to end, or zero for an empty
// interval. Like [time.Time.Sub], the result saturates at the maximum
// [time.Duration] if the interval is too long to represent.
func (p Period) Duration() time.Duration {
	return p.end.Sub(p.start)
}

// Equal reports whether p and other contain exactly the same time points.
// All empty intervals are equal, regardless of their endpoints.
func (p Period) Equal(other Period) bool {
	if p.IsEmpty() && other.IsEmpty() {
		return true
	}

	return p.start.Equal(other.start) && p.end.Equal(other.end)
}

// Contains reports whether t is at or after p's start and strictly before its
// end. Empty intervals contain no time points. The location and monotonic clock
// reading of t do not affect the comparison.
func (p Period) Contains(t time.Time) bool {
	return !t.Before(p.start) && t.Before(p.end)
}

// ContainsPeriod reports whether every time point in other is contained in p.
// Every interval, including an empty one, contains every empty interval.
func (p Period) ContainsPeriod(other Period) bool {
	if other.IsEmpty() {
		return true
	}

	return !p.IsEmpty() && !other.start.Before(p.start) && !other.end.After(p.end)
}

// Overlaps reports whether p and other have a non-empty intersection.
// Adjacent intervals and empty intervals do not overlap.
func (p Period) Overlaps(other Period) bool {
	return !p.IsEmpty() && !other.IsEmpty() && p.start.Before(other.end) && other.start.Before(p.end)
}

// Adjacent reports whether p and other are non-empty and the end of one equals
// the start of the other. Adjacent intervals do not overlap but [Merge] combines
// them into one interval.
func (p Period) Adjacent(other Period) bool {
	return !p.IsEmpty() && !other.IsEmpty() && (p.end.Equal(other.start) || other.end.Equal(p.start))
}

// Intersect returns the non-empty interval shared by p and other.
// If they do not overlap, it returns the zero value and false.
func (p Period) Intersect(other Period) (Period, bool) {
	if !p.Overlaps(other) {
		return Period{}, false
	}
	start, end := p.start, p.end
	if other.start.After(start) {
		start = other.start
	}
	if other.end.Before(end) {
		end = other.end
	}

	return Period{start: start, end: end}, true
}

// Subtract removes every interval in others from p and returns the remaining
// non-empty intervals in ascending start order. The results neither overlap nor
// touch. Empty intervals in others are ignored; duplicates, overlapping cuts,
// and cuts outside p are accepted.
//
// Subtract does not modify or reuse the input slice. It returns nil if nothing
// remains, or a new slice containing p if no cut affects a non-empty p.
// For n cuts, it uses O(n log n) time and O(n) additional space.
func (p Period) Subtract(others ...Period) []Period {
	if p.IsEmpty() {
		return nil
	}
	var result []Period
	cursor := p.start
	for _, cut := range Merge(others...) {
		if !cut.end.After(cursor) {
			continue
		}
		if !cut.start.Before(p.end) {
			break
		}
		if cut.start.After(cursor) {
			result = append(result, Period{start: cursor, end: cut.start})
		}
		cursor = cut.end
		if !cursor.Before(p.end) {
			return result
		}
	}
	result = append(result, Period{start: cursor, end: p.end})

	return result
}
