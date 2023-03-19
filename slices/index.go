package slices

// Index returns the first index in the slice that contains the wanted value, or -1 if not found.
func Index[V comparable](src []V, wanted V) int {
	for i, val := range src {
		if val == wanted {
			return i
		}
	}
	return -1
}

// IndexOf returns the first index in the slice where the selector returned true, or -1 if not found.
func IndexOf[V any](src []V, selector func(V) bool) int {
	for i, val := range src {
		if selector(val) {
			return i
		}
	}
	return -1
}

// IndexOfErr allows the selector to return an error which is then returned to the caller.
func IndexOfErr[V any](src []V, selector func(V) (bool, error)) (int, error) {
	for i, val := range src {
		ok, err := selector(val)
		if err != nil {
			return -1, err
		}
		if ok {
			return i, nil
		}
	}
	return -1, nil
}
