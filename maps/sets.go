package maps

// NewSet creates a new set.
func NewSet[K comparable](keys ...K) map[K]bool {
	return NewSetFromSlice(true, keys...)
}

// NewSetFromSlice creates a set from the given slice,
// with each entry given the same sentinel value.
func NewSetFromSlice[K comparable, V any](value V, keys ...K) map[K]V {
	set := make(map[K]V)
	for _, key := range keys {
		set[key] = value
	}
	return set
}

// Contains returns true if this set contains the given key.
func Contains[K comparable, V any](set map[K]V, key K) bool {
	_, ok := set[key]
	return ok
}

// ContainsAny returns true if this set contains any given key.
func ContainsAny[K comparable, V any](set map[K]V, keys ...K) bool {
	for _, key := range keys {
		if Contains(set, key) {
			return true
		}
	}
	return false
}

// ContainsAll returns true if this set contains every given key.
func ContainsAll[K comparable, V any](set map[K]V, keys ...K) bool {
	for _, key := range keys {
		if !Contains(set, key) {
			return false
		}
	}
	return true
}

// SubsetOf returns true if every key in this set is in the other set.
func SubsetOf[K comparable, V any](this, other map[K]V) bool {
	for key := range this {
		if !Contains(other, key) {
			return false
		}
	}
	return true
}

// AddAll adds multiple keys to this set with the same sentinel value.
func AddAll[K comparable, V any](set map[K]V, value V, keys ...K) {
	for _, key := range keys {
		set[key] = value
	}
}

// Update adds every key from the other set to this set.
func Update[K comparable, V any](this, other map[K]V) {
	for key, value := range other {
		this[key] = value
	}
}

// RemoveAll removes multiple keys from this set.
func RemoveAll[K comparable, V any](set map[K]V, keys ...K) {
	for _, key := range keys {
		delete(set, key)
	}
}

// Discard removes all keys from this set that exist in the other set.
func Discard[K comparable, V any](this, other map[K]V) {
	for key := range other {
		delete(this, key)
	}
}

// Union returns a new set containing all keys from both sets.
func Union[K comparable, V any](this, other map[K]V) map[K]V {
	set := make(map[K]V)
	Update(set, this)
	Update(set, other)
	return set
}

// Intersection returns a new set containing only keys that exist in both sets.
func Intersection[K comparable, V any](this, other map[K]V) map[K]V {
	set := make(map[K]V)
	for key, value := range this {
		if Contains(other, key) {
			set[key] = value
		}
	}
	return set
}

// Difference returns a new set containing all keys in this set that don't exist in the other set.
func Difference[K comparable, V any](this, other map[K]V) map[K]V {
	set := make(map[K]V)
	Update(set, this)
	Discard(set, other)
	return set
}
