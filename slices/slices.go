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

func Filter[S ~[]E, E any](s S, fn func(E) bool) []E {
	var result []E
	for _, item := range s {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}

func Reduce[S ~[]E, E, R any](s S, fn func(R, E) R) R {
	var result R
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

func Concat[S ~[]E, E any](s1, s2 S) S {
	result := make(S, 0, len(s1)+len(s2))
	result = append(result, s1...)
	result = append(result, s2...)
	return result
}

func IsEmpty[S ~[]E, E any](s S) bool {
	return len(s) == 0
}

func IsNotEmpty[S ~[]E, E any](s S) bool {
	return len(s) > 0
}
