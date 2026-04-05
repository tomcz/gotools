package sets

import (
	"cmp"
	"encoding/json"
)

// Set type for encoding to & decoding from JSON.
// Uses ordered keys to allow us to have stable JSON output.
type Set[K cmp.Ordered] map[K]struct{}

// New set creation from keys.
func New[K cmp.Ordered](keys ...K) Set[K] {
	return NewSetWithValue(struct{}{}, keys...)
}

// MarshalJSON implements json.Marshaler
func (s Set[K]) MarshalJSON() ([]byte, error) {
	return json.Marshal(SortedKeys(s))
}

// UnmarshalJSON implements json.Unmarshaler
func (s *Set[K]) UnmarshalJSON(data []byte) error {
	var keys []K
	err := json.Unmarshal(data, &keys)
	if err != nil {
		return err
	}
	*s = New(keys...)
	return nil
}
