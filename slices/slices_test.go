package slices

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type T struct {
	A string
}

func TestMap(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	assert.Equal(t, []int{2, 4, 6, 8, 10}, Map(s1, func(n int) int { return n * 2 }))

	s2 := []int{1, 2, 3, 4, 5}
	assert.Equal(t, []string{"1", "2", "3", "4", "5"}, Map(s2, strconv.Itoa))

	s3 := []string{"1", "2", "3", "4", "5"}
	assert.Equal(t, []T{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}}, Map(s3, func(s string) T { return T{s} }))
}

func TestEach(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	var result []int
	Each(s1, func(n int) { result = append(result, n*2) })
	assert.Equal(t, []int{2, 4, 6, 8, 10}, result)

	s2 := []int{1, 2, 3, 4, 5}
	var result2 []string
	Each(s2, func(n int) { result2 = append(result2, strconv.Itoa(n)) })
	assert.Equal(t, []string{"1", "2", "3", "4", "5"}, result2)

	s3 := []string{"1", "2", "3", "4", "5"}
	var result3 []T
	Each(s3, func(s string) { result3 = append(result3, T{s}) })
}

func TestFilter(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	assert.Equal(t, []int{2, 4}, Filter(s1, func(n int) bool { return n%2 == 0 }))

	s2 := []int{1, 2, 3, 4, 5}
	assert.Equal(t, []int{1, 3, 5}, Filter(s2, func(n int) bool { return n%2 != 0 }))

	s3 := []string{"1", "2", "3", "4", "5"}
	assert.Equal(t, []string{"1", "2", "3"}, Filter(s3, func(s string) bool { return s < "4" }))

	s4 := []T{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}}
	assert.Equal(t, []T{{"1"}, {"2"}, {"3"}}, Filter(s4, func(t T) bool { return t.A < "4" }))
}

func TestReduce(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	assert.Equal(t, 15, Reduce(s1, func(acc, n int) int { return acc + n }))

	s2 := []int{1, 2, 3, 4, 5}
	assert.Equal(t, 0, Reduce(s2, func(acc, n int) int { return acc * n }))

	s3 := []string{"1", "2", "3", "4", "5"}
	assert.Equal(t, "12345", Reduce(s3, func(acc, s string) string { return acc + s }))
}

func TestReverse(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	assert.Equal(t, []int{5, 4, 3, 2, 1}, Reverse(s1))

	s2 := []string{"1", "2", "3", "4", "5"}
	assert.Equal(t, []string{"5", "4", "3", "2", "1"}, Reverse(s2))

	s3 := []T{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}}
	assert.Equal(t, []T{{"5"}, {"4"}, {"3"}, {"2"}, {"1"}}, Reverse(s3))
}

func TestConcat(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := []int{4, 5, 6}
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, Concat(s1, s2))

	s3 := []string{"1", "2", "3"}
	s4 := []string{"4", "5", "6"}
	assert.Equal(t, []string{"1", "2", "3", "4", "5", "6"}, Concat(s3, s4))

	s5 := []T{{"1"}, {"2"}, {"3"}}
	s6 := []T{{"4"}, {"5"}, {"6"}}
	assert.Equal(t, []T{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}, {"6"}}, Concat(s5, s6))
}

func TestIsEmptyAndIsNotEmpty(t *testing.T) {
	s1 := []int{1, 2, 3}
	assert.False(t, IsEmpty(s1))
	assert.True(t, IsNotEmpty(s1))

	var s2 []int
	assert.True(t, IsEmpty(s2))
	assert.False(t, IsNotEmpty(s2))

	s3 := []string{"1", "2", "3"}
	assert.False(t, IsEmpty(s3))
	assert.True(t, IsNotEmpty(s3))

	var s4 []string
	assert.True(t, IsEmpty(s4))
	assert.False(t, IsNotEmpty(s4))

	s5 := []T{{"1"}, {"2"}, {"3"}}
	assert.False(t, IsEmpty(s5))
	assert.True(t, IsNotEmpty(s5))

	var s6 []T
	assert.True(t, IsEmpty(s6))
	assert.False(t, IsNotEmpty(s6))
}

func TestContains(t *testing.T) {
	s1 := []int{1, 2, 3}
	assert.True(t, Contains(s1, 1))
	assert.False(t, Contains(s1, 4))

	s2 := []string{"1", "2", "3"}
	assert.True(t, Contains(s2, "1"))
	assert.False(t, Contains(s2, "4"))

	s3 := []T{{"1"}, {"2"}, {"3"}}
	assert.True(t, Contains(s3, T{"1"}))
	assert.False(t, Contains(s3, T{"4"}))
}

func TestIndexOf(t *testing.T) {
	s1 := []int{1, 2, 3}
	assert.Equal(t, 0, IndexOf(s1, 1))
	assert.Equal(t, -1, IndexOf(s1, 4))

	s2 := []string{"1", "2", "3"}
	assert.Equal(t, 0, IndexOf(s2, "1"))
	assert.Equal(t, -1, IndexOf(s2, "4"))

	s3 := []T{{"1"}, {"2"}, {"3"}}
	assert.Equal(t, 0, IndexOf(s3, T{"1"}))
	assert.Equal(t, -1, IndexOf(s3, T{"4"}))
}

func TestLastIndexOf(t *testing.T) {
	s1 := []int{1, 2, 3, 1}
	assert.Equal(t, 3, LastIndexOf(s1, 1))
	assert.Equal(t, -1, LastIndexOf(s1, 4))

	s2 := []string{"1", "2", "3", "1"}
	assert.Equal(t, 3, LastIndexOf(s2, "1"))
	assert.Equal(t, -1, LastIndexOf(s2, "4"))

	s3 := []T{{"1"}, {"2"}, {"3"}, {"1"}}
	assert.Equal(t, 3, LastIndexOf(s3, T{"1"}))
	assert.Equal(t, -1, LastIndexOf(s3, T{"4"}))
}

func TestUnique(t *testing.T) {
	s1 := []int{1, 2, 3, 1, 2, 3}
	assert.Equal(t, []int{1, 2, 3}, Unique(s1))

	s2 := []string{"1", "2", "3", "1", "2", "3"}
	assert.Equal(t, []string{"1", "2", "3"}, Unique(s2))

	s3 := []T{{"1"}, {"2"}, {"3"}, {"1"}, {"2"}, {"3"}}
	assert.Equal(t, []T{{"1"}, {"2"}, {"3"}}, Unique(s3))
}

func TestDifference(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := []int{2, 3, 4}
	assert.Equal(t, []int{1}, Difference(s1, s2))

	s3 := []string{"1", "2", "3"}
	s4 := []string{"2", "3", "4"}
	assert.Equal(t, []string{"1"}, Difference(s3, s4))

	s5 := []T{{"1"}, {"2"}, {"3"}}
	s6 := []T{{"2"}, {"3"}, {"4"}}
	assert.Equal(t, []T{{"1"}}, Difference(s5, s6))
}

func TestIntersect(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := []int{2, 3, 4}
	assert.Equal(t, []int{2, 3}, Intersect(s1, s2))

	s3 := []string{"1", "2", "3"}
	s4 := []string{"2", "3", "4"}
	assert.Equal(t, []string{"2", "3"}, Intersect(s3, s4))

	s5 := []T{{"1"}, {"2"}, {"3"}}
	s6 := []T{{"2"}, {"3"}, {"4"}}
	assert.Equal(t, []T{{"2"}, {"3"}}, Intersect(s5, s6))
}

func TestOnly(t *testing.T) {
	s1 := []int{1, 2, 3, 4}
	assert.Equal(t, []int{2, 4}, Only(s1, 2, 4))

	s2 := []string{"1", "2", "3", "4"}
	assert.Equal(t, []string{"2", "4"}, Only(s2, "2", "4"))

	s3 := []T{{"1"}, {"2"}, {"3"}, {"4"}}
	assert.Equal(t, []T{{"2"}, {"4"}}, Only(s3, T{"2"}, T{"4"}))
}

func TestWithout(t *testing.T) {
	s1 := []int{1, 2, 3, 4}
	assert.Equal(t, []int{1, 3}, Without(s1, 2, 4))

	s2 := []string{"1", "2", "3", "4"}
	assert.Equal(t, []string{"1", "3"}, Without(s2, "2", "4"))

	s3 := []T{{"1"}, {"2"}, {"3"}, {"4"}}
	assert.Equal(t, []T{{"1"}, {"3"}, {"4"}}, Without(s3, T{"2"}))
}

func TestPartition(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	odd, even := Partition(s1, func(n int) bool { return n%2 != 0 })
	assert.Equal(t, []int{1, 3, 5}, odd)
	assert.Equal(t, []int{2, 4}, even)

	s2 := []string{"1", "2", "3", "4", "5"}
	lessThan3, greaterThan3 := Partition(s2, func(s string) bool { return s < "3" })
	assert.Equal(t, []string{"1", "2"}, lessThan3)
	assert.Equal(t, []string{"3", "4", "5"}, greaterThan3)

	s3 := []T{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}}
	lessThan4, greaterThan4 := Partition(s3, func(t T) bool { return t.A < "4" })
	assert.Equal(t, []T{{"1"}, {"2"}, {"3"}}, lessThan4)
	assert.Equal(t, []T{{"4"}, {"5"}}, greaterThan4)
}

func TestChunk(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5}}, Chunk(s1, 2))

	s2 := []string{"1", "2", "3", "4", "5"}
	assert.Equal(t, [][]string{{"1", "2"}, {"3", "4"}, {"5"}}, Chunk(s2, 2))

	s3 := []T{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}}
	assert.Equal(t, [][]T{{{"1"}, {"2"}}, {{"3"}, {"4"}}, {{"5"}}}, Chunk(s3, 2))
}
