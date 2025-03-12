package utilities

// Map applies the given function to each element in the slice and returns a new slice with the results.
func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}
