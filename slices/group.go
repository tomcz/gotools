package slices

// GroupBy indexes the given values by the result of the key function.
func GroupBy[K comparable, V any](src []V, keyFn func(V) K) map[K][]V {
	dest := make(map[K][]V)
	for _, value := range src {
		key := keyFn(value)
		group := dest[key]
		dest[key] = append(group, value)
	}
	return dest
}
