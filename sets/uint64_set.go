// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package sets

import (
	"encoding/json"
	"sort"
)

// Uint64Set represents a unique collection of uint64 values.
type Uint64Set map[uint64]struct{}

var _ json.Marshaler = Uint64Set{}
var _ json.Unmarshaler = &Uint64Set{}

// NewUint64Set creates a new set of uint64 values.
func NewUint64Set(values ...uint64) Uint64Set {
	s := Uint64Set{}
	s.AddAll(values)
	return s
}

// Contains returns true if this set contains the given value.
func (s Uint64Set) Contains(value uint64) bool {
	_, found := s[value]
	return found
}

// ContainsAny returns true if this set contains any value in the slice.
func (s Uint64Set) ContainsAny(values []uint64) bool {
	for _, value := range values {
		if s.Contains(value) {
			return true
		}
	}
	return false
}

// ContainsAll returns true if this set contains every value in the slice.
func (s Uint64Set) ContainsAll(values []uint64) bool {
	for _, value := range values {
		if !s.Contains(value) {
			return false
		}
	}
	return true
}

// SubsetOf returns true if every element in this set is in the other set.
func (s Uint64Set) SubsetOf(o Uint64Set) bool {
	for value := range s {
		if !o.Contains(value) {
			return false
		}
	}
	return true
}

// Add adds a single value to this set.
func (s Uint64Set) Add(value uint64) {
	s[value] = struct{}{}
}

// AddAll adds multiple values to this set.
func (s Uint64Set) AddAll(values []uint64) {
	for _, value := range values {
		s.Add(value)
	}
}

// Update adds every value from the other set to this set.
func (s Uint64Set) Update(o Uint64Set) {
	for value := range o {
		s.Add(value)
	}
}

// Remove removes a single value from this set.
// If the value is not present this function is a no-op.
func (s Uint64Set) Remove(value uint64) {
	delete(s, value)
}

// RemoveAll removes multiple values from this set.
func (s Uint64Set) RemoveAll(values []uint64) {
	for _, value := range values {
		s.Remove(value)
	}
}

// Discard removes all values in this set that exist in the other set.
func (s Uint64Set) Discard(o Uint64Set) {
	for value := range o {
		s.Remove(value)
	}
}

// Union returns a new set containing all values from this and the other set.
func (s Uint64Set) Union(o Uint64Set) Uint64Set {
	ret := Uint64Set{}
	ret.Update(s)
	ret.Update(o)
	return ret
}

// Intersection returns a new set containing only values that exist in both sets.
func (s Uint64Set) Intersection(o Uint64Set) Uint64Set {
	ret := Uint64Set{}
	for value := range o {
		if s.Contains(value) {
			ret.Add(value)
		}
	}
	return ret
}

// Difference returns a new set containing all values in this set that don't exist in the other set.
func (s Uint64Set) Difference(o Uint64Set) Uint64Set {
	ret := Uint64Set{}
	ret.Update(s)
	ret.Discard(o)
	return ret
}

// Ordered returns an ordered slice of values from this set.
func (s Uint64Set) Ordered() []uint64 {
	ret := make([]uint64, 0, len(s))
	for value := range s {
		ret = append(ret, value)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i] < ret[j]
	})
	return ret
}

// MarshalJSON implements json.Marshaler.
func (s Uint64Set) MarshalJSON() ([]byte, error) {
	var values []uint64
	if s != nil {
		values = s.Ordered()
	}
	return json.Marshal(values)
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *Uint64Set) UnmarshalJSON(in []byte) error {
	var values []uint64
	if err := json.Unmarshal(in, &values); err != nil {
		return err
	}
	*s = Uint64Set{}
	s.AddAll(values)
	return nil
}
