package sets

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	set := NewIntSet(1, 2, 3, 3, 4, 5, 5, 6)

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, set.Ordered())

	assert.True(t, set.Contains(2))
	assert.False(t, set.Contains(10))

	assert.True(t, set.ContainsAny([]int{8, 7, 6}))
	assert.False(t, set.ContainsAny([]int{9, 8, 7}))

	assert.True(t, set.ContainsAll([]int{1, 3, 5}))
	assert.False(t, set.ContainsAll([]int{1, 3, 7}))

	assert.True(t, set.SubsetOf(NewIntSet(1, 2, 3, 4, 5, 6, 7, 8, 9)))
	assert.False(t, set.SubsetOf(NewIntSet(2, 3, 4, 5, 6, 7, 8, 9)))

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, set.Union(NewIntSet(5, 6, 7, 8)).Ordered())
	assert.Equal(t, []int{1, 2, 3, 4}, set.Difference(NewIntSet(5, 6, 7, 8)).Ordered())
	assert.Equal(t, []int{5, 6}, set.Intersection(NewIntSet(5, 6, 7, 8)).Ordered())

	var dst IntSet
	const txt = `[1, 2, 3, 3, 4, 5, 5, 6]`
	err := json.Unmarshal([]byte(txt), &dst)
	if assert.NoError(t, err) {
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, dst.Ordered())
	}

	buf, err := json.Marshal(set)
	if assert.NoError(t, err) {
		assert.Equal(t, `[1,2,3,4,5,6]`, string(buf))
	}
}
