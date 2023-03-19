package slices

// Contains returns true when the slice contains the wanted value.
func Contains[V comparable](src []V, wanted V) bool {
	return Index(src, wanted) >= 0
}

// ContainsAny returns true when the slice contains a value for which the selector returns true.
func ContainsAny[V any](src []V, selector func(V) bool) bool {
	return IndexOf(src, selector) >= 0
}

// ContainsAnyErr allows the selector to return an error which is then returned to the caller.
func ContainsAnyErr[V any](src []V, selector func(V) (bool, error)) (bool, error) {
	index, err := IndexOfErr(src, selector)
	if err != nil {
		return false, err
	}
	return index >= 0, nil
}

// ContainsAll returns true when the selector returns true for all values in the slice.
func ContainsAll[V any](src []V, selector func(V) bool) bool {
	for _, val := range src {
		if !selector(val) {
			return false
		}
	}
	return true
}

// ContainsAllErr allows the selector to return an error which is then returned to the caller.
func ContainsAllErr[V any](src []V, selector func(V) (bool, error)) (bool, error) {
	for _, val := range src {
		ok, err := selector(val)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}
