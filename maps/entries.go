package maps

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// Keys returns a randomly-sorted slice of map keys.
func Keys[K comparable, V any](src map[K]V) []K {
	keys := make([]K, 0, len(src))
	for key := range src {
		keys = append(keys, key)
	}
	return keys
}

// SortedKeys returns a sorted slice of map keys.
func SortedKeys[K constraints.Ordered, V any](src map[K]V) []K {
	keys := Keys(src)
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

// Values returns a randomly-sorted slice of map values.
func Values[K comparable, V any](src map[K]V) []V {
	values := make([]V, 0, len(src))
	for _, value := range src {
		values = append(values, value)
	}
	return values
}

// SortedValues returns a sorted slice of map values.
func SortedValues[K comparable, V constraints.Ordered](src map[K]V) []V {
	values := Values(src)
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})
	return values
}
