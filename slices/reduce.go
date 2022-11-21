package slices

// Reduce a slice to a single value given an initial value and a reducer function.
func Reduce[V any, R any](values []V, acc R, reducer func(val V, acc R) R) R {
	for _, val := range values {
		acc = reducer(val, acc)
	}
	return acc
}

// ReduceErr allows the reducer function to fail and returns the failing error.
func ReduceErr[V any, R any](values []V, acc R, reducer func(val V, acc R) (R, error)) (R, error) {
	var err error
	for _, val := range values {
		acc, err = reducer(val, acc)
		if err != nil {
			return acc, err
		}
	}
	return acc, nil
}
