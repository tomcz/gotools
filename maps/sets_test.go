package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericSets(t *testing.T) {
	set := NewSet(1, 2, 3, 3, 4, 5, 5, 6)

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, SortedKeys(set))

	assert.True(t, Contains(set, 2))
	assert.False(t, Contains(set, 10))

	assert.True(t, ContainsAny(set, 8, 7, 6))
	assert.False(t, ContainsAny(set, 9, 8, 7))

	assert.True(t, ContainsAll(set, 1, 3, 5))
	assert.False(t, ContainsAll(set, 1, 3, 7))

	assert.True(t, SubsetOf(set, NewSet(1, 2, 3, 4, 5, 6, 7, 8, 9)))
	assert.False(t, SubsetOf(set, NewSet(2, 3, 4, 5, 6, 7, 8, 9)))

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, SortedKeys(Union(set, NewSet(5, 6, 7, 8))))
	assert.Equal(t, []int{1, 2, 3, 4}, SortedKeys(Difference(set, NewSet(5, 6, 7, 8))))
	assert.Equal(t, []int{5, 6}, SortedKeys(Intersection(set, NewSet(5, 6, 7, 8))))
}
