package period

import "slices"

// Merge returns the union of periods as non-empty intervals in ascending start
// order. It combines both overlapping and adjacent intervals and ignores empty
// intervals. The results neither overlap nor touch.
//
// Merge does not modify or reuse the input slice and returns nil for empty
// coverage. For n intervals, it uses O(n log n) time and O(n) additional space.
func Merge(periods ...Period) []Period {
	var result []Period
	for _, p := range periods {
		if !p.IsEmpty() {
			result = append(result, p)
		}
	}
	if len(result) == 0 {
		return nil
	}
	slices.SortFunc(result, func(a, b Period) int {
		if order := a.start.Compare(b.start); order != 0 {
			return order
		}

		return a.end.Compare(b.end)
	})
	last := 0
	for _, p := range result[1:] {
		if p.start.After(result[last].end) {
			last++
			result[last] = p
		} else if p.end.After(result[last].end) {
			result[last].end = p.end
		}
	}
	clear(result[last+1:])

	return result[:last+1]
}

// Gaps returns the non-empty gaps between the intervals produced by [Merge],
// in ascending start order. It does not include time before the earliest start
// or after the latest end. Empty intervals do not extend this coverage extent.
// To find gaps within an explicit window, use [Period.Subtract] instead.
//
// Gaps does not modify or reuse the input slice. It returns nil when there are
// fewer than two disjoint merged intervals. For n intervals, it uses O(n log n)
// time and O(n) additional space.
func Gaps(periods ...Period) []Period {
	merged := Merge(periods...)
	if len(merged) < 2 {
		return nil
	}
	result := make([]Period, 0, len(merged)-1)
	for i := 1; i < len(merged); i++ {
		result = append(result, Period{start: merged[i-1].end, end: merged[i].start})
	}

	return result
}
