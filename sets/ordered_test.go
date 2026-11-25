package sets

import (
	"encoding/json"
	"testing"

	"gotest.tools/v3/assert"
)

func TestNew(t *testing.T) {
	set := New(1, 1, 1, 1, 1, 1, 1, 1)
	assert.DeepEqual(t, []int{1}, set.Keys())
}

func TestOrderedSets(t *testing.T) {
	set := New(1, 2, 3, 3, 4, 5, 5, 6)

	assert.DeepEqual(t, []int{1, 2, 3, 4, 5, 6}, set.Keys())

	assert.Assert(t, set.Contains(2))
	assert.Assert(t, !set.Contains(10))

	assert.Assert(t, set.ContainsAny(8, 7, 6))
	assert.Assert(t, !set.ContainsAny(9, 8, 7))

	assert.Assert(t, set.ContainsAll(1, 3, 5))
	assert.Assert(t, !set.ContainsAll(1, 3, 7))

	assert.Assert(t, set.AnyOf(New(4, 5, 6)))
	assert.Assert(t, !set.AnyOf(New(7, 8, 9)))

	assert.Assert(t, set.SubsetOf(New(1, 2, 3, 4, 5, 6, 7, 8, 9)))
	assert.Assert(t, !set.SubsetOf(New(2, 3, 4, 5, 6, 7, 8, 9)))

	set.Add(101)
	set.AddAll(101, 102)
	set.Update(New(102, 103))
	expected := []int{1, 2, 3, 4, 5, 6, 101, 102, 103}
	assert.DeepEqual(t, expected, set.Keys())

	set.Remove(101, 102)
	set.Discard(New(102, 103))
	expected = []int{1, 2, 3, 4, 5, 6}
	assert.DeepEqual(t, expected, set.Keys())

	assert.DeepEqual(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, set.Union(New(5, 6, 7, 8)).Keys())
	assert.DeepEqual(t, []int{1, 2, 3, 4, 7, 8}, set.SymmetricDifference(New(5, 6, 7, 8)).Keys())
	assert.DeepEqual(t, []int{1, 2, 3, 4}, set.Difference(New(5, 6, 7, 8)).Keys())
	assert.DeepEqual(t, []int{5, 6}, set.Intersection(New(5, 6, 7, 8)).Keys())
}

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

	assert.DeepEqual(t, keys, set.Keys())
}
