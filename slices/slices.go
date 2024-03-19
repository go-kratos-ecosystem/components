package slices

func Map[S ~[]E, E, R any](s S, fn func(E) R) []R {
	result := make([]R, 0, len(s))
	for _, item := range s {
		result = append(result, fn(item))
	}
	return result
}

func Each[S ~[]E, E any](s S, fn func(E)) {
	for _, item := range s {
		fn(item)
	}
}

func Prepend[S ~[]E, E any](s S, items ...E) S {
	return append(items, s...)
}

func Append[S ~[]E, E any](s S, items ...E) S {
	return append(s, items...)
}

func Filter[S ~[]E, E any](s S, fn func(E) bool) []E {
	var result []E
	for _, item := range s {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}

func Reduce[S ~[]E, E, R any](s S, fn func(R, E) R, defaults ...R) R {
	var result R
	if len(defaults) > 0 {
		result = defaults[0]
	}
	for _, item := range s {
		result = fn(result, item)
	}
	return result
}

func Reverse[S ~[]E, E any](s S) S {
	result := make(S, len(s))
	for i := 0; i < len(s); i++ {
		result[i] = s[len(s)-1-i]
	}
	return result
}

func Concat[S ~[]E, E any](slices ...S) S {
	var result S
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

func IsEmpty[S ~[]E, E any](s S) bool {
	return len(s) == 0
}

func IsNotEmpty[S ~[]E, E any](s S) bool {
	return len(s) > 0
}

func Contains[S ~[]E, E comparable](s S, e E) bool {
	for _, item := range s {
		if item == e {
			return true
		}
	}
	return false
}

func IndexOf[S ~[]E, E comparable](s S, e E) int {
	for i, item := range s {
		if item == e {
			return i
		}
	}
	return -1
}

func LastIndexOf[S ~[]E, E comparable](s S, e E) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == e {
			return i
		}
	}
	return -1
}

func Unique[S ~[]E, E comparable](s S) S {
	var result S
	seeds := make(map[E]struct{})
	for _, item := range s {
		if _, ok := seeds[item]; !ok {
			seeds[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func UniqueBy[S ~[]E, E any, K comparable](s S, fn func(E) K) S {
	var result S
	seeds := make(map[K]struct{})
	for _, item := range s {
		key := fn(item)
		if _, ok := seeds[key]; !ok {
			seeds[key] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func Difference[S ~[]E, E comparable](s1, s2 S) S {
	var result S
	for _, item := range s1 {
		if !Contains(s2, item) {
			result = append(result, item)
		}
	}
	return result
}

func Intersect[S ~[]E, E comparable](s1, s2 S) S {
	var result S
	for _, item := range s1 {
		if Contains(s2, item) {
			result = append(result, item)
		}
	}
	return result
}

func Only[S ~[]E, E comparable](s S, items ...E) S {
	var result S
	for _, item := range s {
		if Contains(items, item) {
			result = append(result, item)
		}
	}
	return result
}

func Without[S ~[]E, E comparable](s S, items ...E) S {
	var result S
	for _, item := range s {
		if !Contains(items, item) {
			result = append(result, item)
		}
	}
	return result
}

func Partition[S ~[]E, E any](s S, fn func(E) bool) (yes, no S) {
	for _, item := range s {
		if fn(item) {
			yes = append(yes, item)
		} else {
			no = append(no, item)
		}
	}
	return
}

func Chunk[S ~[]E, E any](s S, size int) (result []S) {
	length := len(s)
	for i := 0; i < length; i += size {
		end := i + size
		if end > length {
			end = length
		}
		result = append(result, s[i:end])
	}
	return
}

func GroupBy[S ~[]E, E any, K comparable](s S, fn func(E) K) map[K]S {
	result := make(map[K]S)
	for _, item := range s {
		key := fn(item)
		result[key] = append(result[key], item)
	}
	return result
}

func First[S ~[]E, E any](s S) (E, bool) {
	if len(s) == 0 {
		var zero E
		return zero, false
	}
	return s[0], true
}

func Last[S ~[]E, E any](s S) (E, bool) {
	if len(s) == 0 {
		var zero E
		return zero, false
	}
	return s[len(s)-1], true
}

func Find[S ~[]E, E any](s S, fn func(E) bool) (E, bool) {
	for _, item := range s {
		if fn(item) {
			return item, true
		}
	}
	var zero E
	return zero, false
}

func Index[S ~[]E, E any](s S, fn func(E) bool) (int, bool) {
	for i, item := range s {
		if fn(item) {
			return i, true
		}
	}
	return -1, false
}

func LastIndex[S ~[]E, E any](s S, fn func(E) bool) (int, bool) {
	for i := len(s) - 1; i >= 0; i-- {
		if fn(s[i]) {
			return i, true
		}
	}
	return -1, false
}

func Fill[S ~[]E, E any](s S, value E) S {
	for i := range s {
		s[i] = value
	}
	return s
}
