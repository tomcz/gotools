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

// Entry is a key/value tuple representing a map entry.
type Entry[K comparable, V any] struct {
	Key K
	Val V
}

// Entries returns a randomly-sorted slice of map entries.
func Entries[K comparable, V any](src map[K]V) []Entry[K, V] {
	dest := make([]Entry[K, V], 0, len(src))
	for key, val := range src {
		dest = append(dest, Entry[K, V]{Key: key, Val: val})
	}
	return dest
}

// SortedEntries returns a slice of map entries sorted by their keys.
func SortedEntries[K constraints.Ordered, V any](src map[K]V) []Entry[K, V] {
	entries := Entries(src)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries
}

// FromEntries creates a map from a slice of key/value tuples.
func FromEntries[K comparable, V any](src []Entry[K, V]) map[K]V {
	dest := make(map[K]V)
	for _, e := range src {
		dest[e.Key] = e.Val
	}
	return dest
}
