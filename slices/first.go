package slices

// First returns the first value in the slice where the selector returned true,
// or the zero value of V if the selector never returns true.
func First[V any](src []V, selector func(V) bool) V {
	for _, val := range src {
		if selector(val) {
			return val
		}
	}
	var zero V
	return zero
}

// FirstErr allows the selector to return an error which is then returned to the caller.
func FirstErr[V any](src []V, selector func(V) (bool, error)) (V, error) {
	var zero V
	for _, val := range src {
		ok, err := selector(val)
		if err != nil {
			return zero, err
		}
		if ok {
			return val, nil
		}
	}
	return zero, nil
}
