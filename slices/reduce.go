package slices

// Reduce a slice to a single value given an initial value.
func Reduce[V any, R any](values []V, acc R, fn func(val V, acc R) R) R {
	for _, val := range values {
		acc = fn(val, acc)
	}
	return acc
}

// ReduceErr allows the reduction to fail for any value that returns an error.
func ReduceErr[V any, R any](values []V, acc R, fn func(val V, acc R) (R, error)) (R, error) {
	var err error
	for _, val := range values {
		acc, err = fn(val, acc)
		if err != nil {
			return acc, err
		}
	}
	return acc, nil
}
