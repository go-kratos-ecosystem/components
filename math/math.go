package math

import (
	"errors"

	"github.com/go-kratos-ecosystem/components/v2/constraints"
)

var ErrorEmpty = errors.New("math: the input list is empty")

func Max[T constraints.Ordered](nums ...T) (T, error) {
	if len(nums) == 0 {
		var zero T
		return zero, ErrorEmpty
	}

	m := nums[0]
	for _, num := range nums {
		if num > m {
			m = num
		}
	}
	return m, nil
}

func Min[T constraints.Ordered](nums ...T) (T, error) {
	if len(nums) == 0 {
		var zero T
		return zero, ErrorEmpty
	}

	m := nums[0]
	for _, num := range nums {
		if num < m {
			m = num
		}
	}

	return m, nil
}

func Sum[T constraints.Ordered](nums ...T) (T, error) {
	if len(nums) == 0 {
		var zero T
		return zero, ErrorEmpty
	}

	var sum T
	for _, num := range nums {
		sum += num
	}
	return sum, nil
}

func Average[T constraints.Integer | constraints.Float](nums ...T) (T, error) {
	if len(nums) == 0 {
		var zero T
		return zero, ErrorEmpty
	}

	sum, err := Sum(nums...)
	if err != nil {
		return sum, err
	}
	return sum / T(len(nums)), nil
}
