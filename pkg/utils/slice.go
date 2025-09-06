package utils

// ToAnySlice converts a slice of any type to a slice of empty interface.
// It takes a slice of type T and returns a slice of type []any, allowing
// for heterogeneous collections of values.
func ToAnySlice[T any](slice []T) []any {
	result := make([]any, len(slice))
	for i := range slice {
		result[i] = slice[i]
	}
	return result
}

// TransformSlice applies a transformation function to each element of a slice.
// It takes a slice of type T and a transform function that converts type T to type K,
// returning a new slice of type []K with all transformed elements.
func TransformSlice[T, K any](slice []T, transform func(T) K) []K {
	result := make([]K, len(slice))
	for i, v := range slice {
		result[i] = transform(v)
	}
	return result
}
