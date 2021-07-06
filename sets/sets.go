package sets

import (
	"encoding/json"
	"sort"

	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=sets_gen.go gen "Something=string,int,int64,uint64"

type Something generic.Number

// SomethingSet adds set semantics to a map of values.
type SomethingSet map[Something]bool

var _ json.Marshaler = SomethingSet{}
var _ json.Unmarshaler = &SomethingSet{}

// NewSomethingSet creates a new set of Something.
func NewSomethingSet(values ...Something) SomethingSet {
	s := SomethingSet{}
	s.AddAll(values)
	return s
}

// Contains returns true if this set contains the given value.
func (s SomethingSet) Contains(value Something) bool {
	return s[value]
}

// ContainsAny returns true if this set contains any value in the slice.
func (s SomethingSet) ContainsAny(values []Something) bool {
	for _, value := range values {
		if s.Contains(value) {
			return true
		}
	}
	return false
}

// ContainsAll returns true if this set contains every value in the slice.
func (s SomethingSet) ContainsAll(values []Something) bool {
	for _, value := range values {
		if !s.Contains(value) {
			return false
		}
	}
	return true
}

// SubsetOf returns true if every element in this set is in the other set.
func (s SomethingSet) SubsetOf(o SomethingSet) bool {
	for value := range s {
		if !o.Contains(value) {
			return false
		}
	}
	return true
}

// Add adds a single value to this set.
func (s SomethingSet) Add(value Something) {
	s[value] = true
}

// AddAll adds multiple values to this set.
func (s SomethingSet) AddAll(values []Something) {
	for _, value := range values {
		s.Add(value)
	}
}

// Update adds every value from the other set to this set.
func (s SomethingSet) Update(o SomethingSet) {
	for value := range o {
		s.Add(value)
	}
}

// Remove removes a single value from this set.
// If the value is not present this function is a no-op.
func (s SomethingSet) Remove(value Something) {
	delete(s, value)
}

// RemoveAll removes multiple values from this set.
func (s SomethingSet) RemoveAll(values []Something) {
	for _, value := range values {
		s.Remove(value)
	}
}

// Discard removes all values in this set that exist in the other set.
func (s SomethingSet) Discard(o SomethingSet) {
	for value := range o {
		s.Remove(value)
	}
}

// Union returns a new set containing all values from this and the other set.
func (s SomethingSet) Union(o SomethingSet) SomethingSet {
	ret := SomethingSet{}
	ret.Update(s)
	ret.Update(o)
	return ret
}

// Intersection returns a new set containing only values that exist in both sets.
func (s SomethingSet) Intersection(o SomethingSet) SomethingSet {
	ret := SomethingSet{}
	for value := range o {
		if s.Contains(value) {
			ret.Add(value)
		}
	}
	return ret
}

// Difference returns a new set containing all values in this set that don't exist in the other set.
func (s SomethingSet) Difference(o SomethingSet) SomethingSet {
	ret := SomethingSet{}
	ret.Update(s)
	ret.Discard(o)
	return ret
}

// Ordered returns an ordered slice of values from this set.
func (s SomethingSet) Ordered() []Something {
	ret := make([]Something, 0, len(s))
	for value := range s {
		ret = append(ret, value)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i] < ret[j]
	})
	return ret
}

func (s SomethingSet) MarshalJSON() ([]byte, error) {
	var values []Something
	if s != nil {
		values = s.Ordered()
	}
	return json.Marshal(values)
}

func (s *SomethingSet) UnmarshalJSON(in []byte) error {
	var values []Something
	if err := json.Unmarshal(in, &values); err != nil {
		return err
	}
	*s = SomethingSet{}
	s.AddAll(values)
	return nil
}
