package period_test

import (
	"slices"
	"testing"
	"time"

	"github.com/go-fries/fries/period/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMerge(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input [][2]int
		want  [][2]int
	}{
		{name: "no intervals"},
		{name: "empty intervals", input: [][2]int{{1, 1}, {2, 2}}},
		{name: "single", input: [][2]int{{1, 3}}, want: [][2]int{{1, 3}}},
		{name: "adjacent", input: [][2]int{{1, 2}, {2, 3}}, want: [][2]int{{1, 3}}},
		{name: "overlapping", input: [][2]int{{1, 3}, {2, 4}}, want: [][2]int{{1, 4}}},
		{name: "unsorted disjoint", input: [][2]int{{5, 6}, {1, 2}, {3, 4}}, want: [][2]int{{1, 2}, {3, 4}, {5, 6}}},
		{name: "same start and duplicates", input: [][2]int{{1, 5}, {1, 2}, {1, 5}, {1, 3}}, want: [][2]int{{1, 5}}},
		{name: "nested", input: [][2]int{{1, 10}, {2, 4}, {5, 7}}, want: [][2]int{{1, 10}}},
		{name: "empty does not bridge gap", input: [][2]int{{1, 2}, {3, 3}, {4, 5}}, want: [][2]int{{1, 2}, {4, 5}}},
		{name: "chain", input: [][2]int{{8, 10}, {1, 3}, {2, 5}, {7, 9}, {5, 7}}, want: [][2]int{{1, 10}}},
		{name: "multiple groups", input: [][2]int{{5, 8}, {1, 3}, {6, 7}, {2, 4}, {10, 12}, {9, 11}}, want: [][2]int{{1, 4}, {5, 8}, {9, 12}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			input := intervals(t, tt.input...)
			got := period.Merge(input...)
			assert.Equal(t, intervals(t, tt.want...), got)
			assert.Equal(t, intervals(t, tt.input...), input)
			if len(got) > 0 {
				got[0] = period.Period{}
				assert.Equal(t, intervals(t, tt.input...), input)
			}
		})
	}
	assert.Nil(t, period.Merge(period.Period{}))
}

func TestGaps(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input [][2]int
		want  [][2]int
	}{
		{name: "no intervals"},
		{name: "empty intervals", input: [][2]int{{0, 0}, {10, 10}}},
		{name: "single", input: [][2]int{{1, 5}}},
		{name: "adjacent and overlapping", input: [][2]int{{3, 6}, {1, 3}, {5, 8}}},
		{name: "unsorted", input: [][2]int{{7, 9}, {1, 3}, {5, 6}}, want: [][2]int{{3, 5}, {6, 7}}},
		{name: "merge before finding gaps", input: [][2]int{{1, 4}, {2, 3}, {8, 9}, {3, 6}, {9, 10}}, want: [][2]int{{6, 8}}},
		{name: "empty endpoints ignored", input: [][2]int{{0, 0}, {1, 3}, {5, 7}, {10, 10}}, want: [][2]int{{3, 5}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			input := intervals(t, tt.input...)
			got := period.Gaps(input...)
			assert.Equal(t, intervals(t, tt.want...), got)
			if len(got) > 0 {
				got[0] = period.Period{}
			}
			assert.Equal(t, intervals(t, tt.input...), input)
		})
	}
}

func TestSetOperationsAcrossLocations(t *testing.T) {
	start := time.Date(2026, time.September, 12, 9, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	first, err := period.New(start, start.Add(time.Hour))
	require.NoError(t, err)
	second, err := period.New(start.Add(time.Hour).UTC(), start.Add(2*time.Hour).UTC())
	require.NoError(t, err)
	duplicate, err := period.New(start.UTC(), start.Add(time.Hour).UTC())
	require.NoError(t, err)
	assert.True(t, first.Equal(duplicate))
	assert.True(t, first.Adjacent(second))
	merged := period.Merge(second, first, duplicate)
	require.Len(t, merged, 1)
	assert.Equal(t, 2*time.Hour, merged[0].Duration())
	assert.Equal(t, []period.Period{second}, merged[0].Subtract(first))
	assert.Nil(t, period.Gaps(second, first))
}

// FuzzSetOperations checks membership against a discrete nanosecond bitmap,
// independently of the interval sorting and cursor algorithms.
func FuzzSetOperations(f *testing.F) {
	f.Add([]byte{0, 31})
	f.Add([]byte{0, 31, 2, 5, 4, 8, 10, 12, 12, 20, 25, 30})
	f.Add([]byte{5, 20, 0, 31})
	f.Add([]byte{3, 3, 1, 1, 2, 2})
	f.Add([]byte{0, 10, 7, 9, 1, 2, 3, 5})
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) < 2 {
			return
		}
		data = data[:min(len(data), 66)]
		bounds := func(a, b byte) (int, int) {
			start, end := int(a%32), int(b%32)

			return min(start, end), max(start, end)
		}
		start, end := bounds(data[0], data[1])
		base := interval(t, start, end)
		var cuts []period.Period
		var covered [32]bool
		for i := 2; i+1 < len(data); i += 2 {
			a, b := bounds(data[i], data[i+1])
			cuts = append(cuts, interval(t, a, b))
			for n := a; n < b; n++ {
				covered[n] = true
			}
		}
		before := slices.Clone(cuts)
		merged := period.Merge(cuts...)
		remaining := base.Subtract(cuts...)
		gaps := period.Gaps(cuts...)
		require.Equal(t, before, cuts)
		require.Equal(t, merged, period.Merge(merged...))
		require.Equal(t, remaining, period.Merge(remaining...))
		require.Equal(t, gaps, period.Merge(gaps...))
		slices.Reverse(cuts)
		require.Equal(t, merged, period.Merge(cuts...))
		first, last := 32, -1
		for n, contains := range covered {
			if contains {
				first = min(first, n)
				last = n
			}
		}
		for n := -1; n <= 32; n++ {
			wantCovered := n >= 0 && n < 32 && covered[n]
			wantRemaining := n >= start && n < end && !wantCovered
			wantGap := n >= first && n <= last && !wantCovered
			assert.Equal(t, wantCovered, containsPoint(merged, instant(n)), "union at %d", n)
			assert.Equal(t, wantRemaining, containsPoint(remaining, instant(n)), "difference at %d", n)
			assert.Equal(t, wantGap, containsPoint(gaps, instant(n)), "gap at %d", n)
		}
		for _, result := range [][]period.Period{merged, remaining, gaps} {
			for i, p := range result {
				require.False(t, p.IsEmpty())
				if i > 0 {
					require.True(t, result[i-1].End().Before(p.Start()))
				}
			}
		}
	})
}

func containsPoint(periods []period.Period, at time.Time) bool {
	return slices.ContainsFunc(periods, func(p period.Period) bool { return p.Contains(at) })
}
