package sets

import (
	"cmp"
	"encoding/json"
)

// Set type using [cmp.Ordered] keys for predictable JSON output.
type Set[K cmp.Ordered] map[K]struct{}

// New creates an ordered set containing the given keys.
func New[K cmp.Ordered](keys ...K) Set[K] {
	return NewSet(keys...)
}

// Contains returns true if this set contains the given key.
func (s Set[K]) Contains(key K) bool {
	return Contains(s, key)
}

// ContainsAny returns true if this set contains any given key.
func (s Set[K]) ContainsAny(keys ...K) bool {
	return ContainsAny(s, keys...)
}

// ContainsAll returns true if this set contains every given key.
func (s Set[K]) ContainsAll(keys ...K) bool {
	return ContainsAll(s, keys...)
}

// AnyOf returns true if any key in this set is in the other set.
func (s Set[K]) AnyOf(o Set[K]) bool {
	return AnyOf(s, o)
}

// SubsetOf returns true if every key in this set is in the other set.
func (s Set[K]) SubsetOf(o Set[K]) bool {
	return SubsetOf(s, o)
}

// Add a key to this set.
func (s Set[K]) Add(key K) {
	s[key] = empty
}

// AddAll keys to this set.
func (s Set[K]) AddAll(keys ...K) {
	Add(s, keys...)
}

// Update adds every key from the other set to this set.
func (s Set[K]) Update(o Set[K]) {
	Update(s, o)
}

// Remove all keys from this set that exist in the given slice.
func (s Set[K]) Remove(keys ...K) {
	Remove(s, keys...)
}

// Discard all keys from this set that exist in the other set.
func (s Set[K]) Discard(o Set[K]) {
	Discard(s, o)
}

// Union returns a new set containing all keys from both sets.
func (s Set[K]) Union(o Set[K]) Set[K] {
	return Union(s, o)
}

// Intersection returns a new set containing only keys that exist in both sets.
func (s Set[K]) Intersection(o Set[K]) Set[K] {
	return Intersection(s, o)
}

// Difference returns a new set containing all keys in this set that don't exist in the other set.
func (s Set[K]) Difference(o Set[K]) Set[K] {
	return Difference(s, o)
}

// SymmetricDifference returns a new set with keys in either this set or other but not both.
func (s Set[K]) SymmetricDifference(o Set[K]) Set[K] {
	return SymmetricDifference(s, o)
}

// Keys returns a sorted slice of all keys from this set.
func (s Set[K]) Keys() []K {
	return SortedKeys(s)
}

// MarshalJSON implements [json.Marshaler]
func (s Set[K]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Keys())
}

// UnmarshalJSON implements [json.Unmarshaler]
func (s *Set[K]) UnmarshalJSON(data []byte) error {
	var keys []K
	err := json.Unmarshal(data, &keys)
	if err != nil {
		return err
	}
	*s = New(keys...)
	return nil
}
