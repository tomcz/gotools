package sets

import (
	"encoding/json"
	"testing"

	"gotest.tools/v3/assert"
)

func TestSet_MarshalJSON(t *testing.T) {
	keys := []int{1, 2, 3, 4, 5}
	keysJSON, err := json.Marshal(keys)
	assert.NilError(t, err)

	set := New(keys...)
	setJSON, err := json.Marshal(set)
	assert.NilError(t, err)

	assert.Equal(t, string(keysJSON), string(setJSON))
}

func TestSet_UnmarshalJSON(t *testing.T) {
	keys := []int{1, 2, 3, 4, 5}
	keysJSON, err := json.Marshal(keys)
	assert.NilError(t, err)

	var set Set[int]
	err = json.Unmarshal(keysJSON, &set)
	assert.NilError(t, err)

	assert.DeepEqual(t, keys, SortedKeys(set))
}
