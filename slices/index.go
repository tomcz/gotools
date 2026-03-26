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

// IndexBy indexes the given values by the result of the key function.
// Duplicate keys produced by the key function are handled as last-one-wins.
// If you intend for the key function to produce duplicate keys then
// you really should be using the GroupBy function instead.
func IndexBy[K comparable, V any](src []V, keyFn func(V) K) map[K]V {
	dest := make(map[K]V)
	for _, value := range src {
		key := keyFn(value)
		dest[key] = value
	}
	return dest
}

// IndexByErr allows the key function to fail and returns the failing error.
// Duplicate keys produced by the key function are handled as last-one-wins.
// If you intend for the key function to produce duplicate keys then
// you really should be using the GroupByErr function instead.
func IndexByErr[K comparable, V any](src []V, keyFn func(V) (K, error)) (map[K]V, error) {
	dest := make(map[K]V)
	for _, value := range src {
		key, err := keyFn(value)
		if err != nil {
			return nil, err
		}
		dest[key] = value
	}
	return dest, nil
}
