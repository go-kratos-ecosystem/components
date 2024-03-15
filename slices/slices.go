package slices

func Map[S ~[]E, E, R any](s S, fn func(E) R) []R {
	result := make([]R, 0, len(s))
	for _, item := range s {
		result = append(result, fn(item))
	}
	return result
}
