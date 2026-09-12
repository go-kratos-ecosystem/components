package period_test

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/go-fries/fries/time/period/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("normalizes endpoints", func(t *testing.T) {
		start := time.Date(2026, time.September, 12, 9, 0, 0, 1, time.FixedZone("UTC+8", 8*60*60))
		end := time.Date(2026, time.September, 12, 2, 0, 0, 2, time.UTC)
		p, err := period.New(start, end)
		require.NoError(t, err)
		assert.Equal(t, start.UTC(), p.Start())
		assert.Equal(t, end, p.End())
		assert.Same(t, time.UTC, p.Start().Location())
		assert.Equal(t, time.Hour+time.Nanosecond, p.Duration())
	})
	t.Run("removes monotonic clock readings", func(t *testing.T) {
		start := time.Now()
		end := start.Add(time.Second)
		p, err := period.New(start, end)
		require.NoError(t, err)
		wall, err := period.New(start.Round(0).UTC(), end.Round(0).UTC())
		require.NoError(t, err)
		assert.Equal(t, wall, p)
		assert.Equal(t, start.Round(0).UTC(), p.Start())
		assert.True(t, p.Contains(start))
		assert.False(t, p.Contains(end))
	})
	t.Run("equal instants in different locations are empty", func(t *testing.T) {
		start := instant(10)
		end := start.In(time.FixedZone("UTC-5", -5*60*60))
		p, err := period.New(start, end)
		require.NoError(t, err)
		assert.True(t, p.IsEmpty())
		assert.Zero(t, p.Duration())
		assert.False(t, p.Contains(start))
	})
	t.Run("reversed endpoints", func(t *testing.T) {
		p, err := period.New(instant(2), instant(1))
		require.ErrorIs(t, err, period.ErrInvalidRange)
		assert.Equal(t, period.Period{}, p)
	})
	t.Run("zero value", func(t *testing.T) {
		var p period.Period
		assert.True(t, p.IsEmpty())
		assert.True(t, p.Start().IsZero())
		assert.True(t, p.End().IsZero())
		assert.Zero(t, p.Duration())
		assert.False(t, p.Contains(time.Time{}))
		assert.Nil(t, p.Subtract())
	})
}

func TestPeriodContains(t *testing.T) {
	p := interval(t, 1, 3)
	for _, tt := range []struct {
		name string
		at   int
		want bool
	}{
		{name: "before", at: 0},
		{name: "start", at: 1, want: true},
		{name: "inside", at: 2, want: true},
		{name: "end", at: 3},
		{name: "after", at: 4},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, p.Contains(instant(tt.at)))
			assert.Equal(t, tt.want, p.Contains(instant(tt.at).In(time.FixedZone("other", 3600))))
		})
	}
}

func TestPeriodRelations(t *testing.T) {
	p := interval(t, 2, 6)
	for _, tt := range []struct {
		name         string
		start, end   int
		contains     bool
		equal        bool
		overlaps     bool
		adjacent     bool
		intersection [2]int
	}{
		{name: "before", start: 0, end: 1},
		{name: "adjacent before", start: 0, end: 2, adjacent: true},
		{name: "left overlap", start: 1, end: 3, overlaps: true, intersection: [2]int{2, 3}},
		{name: "contained", start: 3, end: 5, contains: true, overlaps: true, intersection: [2]int{3, 5}},
		{name: "shared start", start: 2, end: 4, contains: true, overlaps: true, intersection: [2]int{2, 4}},
		{name: "shared end", start: 4, end: 6, contains: true, overlaps: true, intersection: [2]int{4, 6}},
		{name: "same", start: 2, end: 6, contains: true, equal: true, overlaps: true, intersection: [2]int{2, 6}},
		{name: "enclosing", start: 1, end: 7, overlaps: true, intersection: [2]int{2, 6}},
		{name: "right overlap", start: 5, end: 7, overlaps: true, intersection: [2]int{5, 6}},
		{name: "adjacent after", start: 6, end: 8, adjacent: true},
		{name: "after", start: 7, end: 8},
		{name: "empty before", start: 0, end: 0, contains: true},
		{name: "empty start", start: 2, end: 2, contains: true},
		{name: "empty inside", start: 4, end: 4, contains: true},
		{name: "empty end", start: 6, end: 6, contains: true},
		{name: "empty after", start: 8, end: 8, contains: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			other := interval(t, tt.start, tt.end)
			assert.Equal(t, tt.contains, p.ContainsPeriod(other))
			assert.Equal(t, tt.equal, p.Equal(other))
			assert.Equal(t, tt.equal, other.Equal(p))
			assert.Equal(t, tt.overlaps, p.Overlaps(other))
			assert.Equal(t, tt.overlaps, other.Overlaps(p))
			assert.Equal(t, tt.adjacent, p.Adjacent(other))
			assert.Equal(t, tt.adjacent, other.Adjacent(p))
			intersection, ok := p.Intersect(other)
			assert.Equal(t, tt.overlaps, ok)
			reverse, reverseOK := other.Intersect(p)
			assert.Equal(t, ok, reverseOK)
			assert.Equal(t, intersection, reverse)
			if ok {
				assert.Equal(t, interval(t, tt.intersection[0], tt.intersection[1]), intersection)
			} else {
				assert.Equal(t, period.Period{}, intersection)
			}
		})
	}
}

func TestPeriodEmptySetRelations(t *testing.T) {
	var zero period.Period
	empty := interval(t, 10, 10)
	p := interval(t, 1, 2)
	assert.True(t, zero.Equal(empty))
	assert.True(t, empty.Equal(zero))
	assert.True(t, zero.ContainsPeriod(empty))
	assert.True(t, empty.ContainsPeriod(zero))
	assert.False(t, empty.ContainsPeriod(p))
	assert.False(t, empty.Equal(p))
	assert.False(t, zero.Adjacent(empty))
	assert.False(t, zero.Overlaps(empty))
	intersection, ok := zero.Intersect(empty)
	assert.False(t, ok)
	assert.Equal(t, zero, intersection)
}

func TestPeriodSubtract(t *testing.T) {
	p := interval(t, 2, 10)
	for _, tt := range []struct {
		name string
		cuts [][2]int
		want [][2]int
	}{
		{name: "no cuts", want: [][2]int{{2, 10}}},
		{name: "empty cuts", cuts: [][2]int{{1, 1}, {5, 5}, {12, 12}}, want: [][2]int{{2, 10}}},
		{name: "before and after", cuts: [][2]int{{0, 1}, {11, 12}}, want: [][2]int{{2, 10}}},
		{name: "adjacent", cuts: [][2]int{{0, 2}, {10, 12}}, want: [][2]int{{2, 10}}},
		{name: "middle", cuts: [][2]int{{4, 6}}, want: [][2]int{{2, 4}, {6, 10}}},
		{name: "left edge", cuts: [][2]int{{0, 4}}, want: [][2]int{{4, 10}}},
		{name: "right edge", cuts: [][2]int{{8, 12}}, want: [][2]int{{2, 8}}},
		{name: "same", cuts: [][2]int{{2, 10}}},
		{name: "enclosing", cuts: [][2]int{{0, 12}}},
		{name: "combined coverage", cuts: [][2]int{{6, 12}, {0, 6}}},
		{name: "unordered overlapping and duplicate", cuts: [][2]int{{7, 8}, {3, 5}, {4, 6}, {3, 5}}, want: [][2]int{{2, 3}, {6, 7}, {8, 10}}},
		{name: "multiple separated cuts", cuts: [][2]int{{8, 9}, {3, 4}, {5, 6}}, want: [][2]int{{2, 3}, {4, 5}, {6, 8}, {9, 10}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cuts := intervals(t, tt.cuts...)
			got := p.Subtract(cuts...)
			assert.Equal(t, intervals(t, tt.want...), got)
			assert.Equal(t, intervals(t, tt.cuts...), cuts)
			assert.Equal(t, interval(t, 2, 10), p)
		})
	}
	assert.Nil(t, interval(t, 2, 2).Subtract(p))
}

func TestSubtractResultDoesNotAliasCuts(t *testing.T) {
	p := interval(t, 0, 10)
	cuts := intervals(t, [2]int{4, 6}, [2]int{1, 2})
	before := slices.Clone(cuts)
	got := p.Subtract(cuts...)
	require.NotEmpty(t, got)
	got[0] = period.Period{}
	assert.Equal(t, before, cuts)
}

func TestPeriodDurationAcrossDST(t *testing.T) {
	location, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	for _, tt := range []struct {
		name  string
		month time.Month
		day   int
		want  time.Duration
	}{
		{name: "spring", month: time.March, day: 8, want: 23 * time.Hour},
		{name: "autumn", month: time.November, day: 1, want: 25 * time.Hour},
	} {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Date(2026, tt.month, tt.day, 0, 0, 0, 0, location)
			end := start.AddDate(0, 0, 1)
			p, err := period.New(start, end)
			require.NoError(t, err)
			assert.Equal(t, tt.want, p.Duration())
			assert.True(t, p.Contains(start))
			assert.False(t, p.Contains(end))
		})
	}
}

func TestPeriodLargeRange(t *testing.T) {
	start := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	middle := time.Date(5000, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC)
	p, err := period.New(start, end)
	require.NoError(t, err)
	assert.Equal(t, time.Duration(math.MaxInt64), p.Duration())
	assert.True(t, p.Contains(middle))
	cut, err := period.New(start, middle)
	require.NoError(t, err)
	remaining := p.Subtract(cut)
	require.Len(t, remaining, 1)
	assert.Equal(t, middle, remaining[0].Start())
	assert.Equal(t, end, remaining[0].End())
	assert.Equal(t, []period.Period{p}, period.Merge(cut, remaining[0]))
}

func instant(offset int) time.Time {
	return time.Date(2026, time.September, 12, 0, 0, 0, offset, time.UTC)
}

func interval(t *testing.T, start, end int) period.Period {
	t.Helper()
	p, err := period.New(instant(start), instant(end))
	require.NoError(t, err)

	return p
}

func intervals(t *testing.T, endpoints ...[2]int) []period.Period {
	t.Helper()
	if len(endpoints) == 0 {
		return nil
	}
	result := make([]period.Period, 0, len(endpoints))
	for _, pair := range endpoints {
		result = append(result, interval(t, pair[0], pair[1]))
	}

	return result
}
