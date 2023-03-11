package slices

// Reduce a slice to a single value given an initial value and a reducer function.
func Reduce[V any, A any](values []V, acc A, reducer func(acc A, val V) A) A {
	for _, val := range values {
		acc = reducer(acc, val)
	}
	return acc
}

// ReduceErr allows the reducer function to fail and returns the failing error.
func ReduceErr[V any, A any](values []V, acc A, reducer func(acc A, val V) (A, error)) (A, error) {
	var err error
	for _, val := range values {
		acc, err = reducer(acc, val)
		if err != nil {
			return acc, err
		}
	}
	return acc, nil
}
