package slices

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	assert.Equal(t, []int{2, 4, 6, 8, 10}, Map(s1, func(n int) int { return n * 2 }))

	s2 := []int{1, 2, 3, 4, 5}
	assert.Equal(t, []string{"1", "2", "3", "4", "5"}, Map(s2, strconv.Itoa))

	s3 := []string{"1", "2", "3", "4", "5"}
	type T struct {
		A string
	}
	assert.Equal(t, []T{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}}, Map(s3, func(s string) T { return T{s} }))
}
