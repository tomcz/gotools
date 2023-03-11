package maps

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestGenericSets(t *testing.T) {
	set := NewSet(1, 2, 3, 3, 4, 5, 5, 6)

	assert.DeepEqual(t, []int{1, 2, 3, 4, 5, 6}, SortedKeys(set))

	assert.Assert(t, Contains(set, 2))
	assert.Assert(t, !Contains(set, 10))

	assert.Assert(t, ContainsAny(set, 8, 7, 6))
	assert.Assert(t, !ContainsAny(set, 9, 8, 7))

	assert.Assert(t, ContainsAll(set, 1, 3, 5))
	assert.Assert(t, !ContainsAll(set, 1, 3, 7))

	assert.Assert(t, SubsetOf(set, NewSet(1, 2, 3, 4, 5, 6, 7, 8, 9)))
	assert.Assert(t, !SubsetOf(set, NewSet(2, 3, 4, 5, 6, 7, 8, 9)))

	AddAll(set, true, 101, 102, 103)
	expected := []int{1, 2, 3, 4, 5, 6, 101, 102, 103}
	assert.DeepEqual(t, expected, SortedKeys(set))

	RemoveAll(set, 101, 102, 103)
	expected = []int{1, 2, 3, 4, 5, 6}
	assert.DeepEqual(t, expected, SortedKeys(set))

	assert.DeepEqual(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, SortedKeys(Union(set, NewSet(5, 6, 7, 8))))
	assert.DeepEqual(t, []int{1, 2, 3, 4}, SortedKeys(Difference(set, NewSet(5, 6, 7, 8))))
	assert.DeepEqual(t, []int{5, 6}, SortedKeys(Intersection(set, NewSet(5, 6, 7, 8))))
}
