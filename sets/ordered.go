package sets

import (
	"cmp"
	"encoding/json"
)

// Set type for encoding to & decoding from JSON.
// Uses [cmp.Ordered] keys to allow us to have stable JSON output.
type Set[K cmp.Ordered] map[K]bool

// New creates an ordered set using [NewSet].
func New[K cmp.Ordered](keys ...K) Set[K] {
	return NewSet(keys...)
}

// MarshalJSON implements [json.Marshaler]
func (s Set[K]) MarshalJSON() ([]byte, error) {
	return json.Marshal(SortedKeys(s))
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
