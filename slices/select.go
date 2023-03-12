package slices

// Select returns a slice with all values where the selector returned true.
func Select[V any](src []V, selector func(V) bool) []V {
	var dest []V
	for _, val := range src {
		if selector(val) {
			dest = append(dest, val)
		}
	}
	return dest
}

// SelectErr allows the selector to return an error which is then returned to the caller.
func SelectErr[V any](src []V, selector func(V) (bool, error)) ([]V, error) {
	var dest []V
	for _, val := range src {
		ok, err := selector(val)
		if err != nil {
			return nil, err
		}
		if ok {
			dest = append(dest, val)
		}
	}
	return dest, nil
}
