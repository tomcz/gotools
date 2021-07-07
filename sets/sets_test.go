package sets

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	set := NewSomethingSet(1, 2, 3, 3, 4, 5, 5, 6)

	assert.Equal(t, []Something{1, 2, 3, 4, 5, 6}, set.Ordered())

	assert.True(t, set.Contains(2))
	assert.False(t, set.Contains(10))

	assert.True(t, set.ContainsAny([]Something{8, 7, 6}))
	assert.False(t, set.ContainsAny([]Something{9, 8, 7}))

	assert.True(t, set.ContainsAll([]Something{1, 3, 5}))
	assert.False(t, set.ContainsAll([]Something{1, 3, 7}))

	assert.True(t, set.SubsetOf(NewSomethingSet(1, 2, 3, 4, 5, 6, 7, 8, 9)))
	assert.False(t, set.SubsetOf(NewSomethingSet(2, 3, 4, 5, 6, 7, 8, 9)))

	assert.Equal(t, []Something{1, 2, 3, 4, 5, 6, 7, 8}, set.Union(NewSomethingSet(5, 6, 7, 8)).Ordered())
	assert.Equal(t, []Something{1, 2, 3, 4}, set.Difference(NewSomethingSet(5, 6, 7, 8)).Ordered())
	assert.Equal(t, []Something{5, 6}, set.Intersection(NewSomethingSet(5, 6, 7, 8)).Ordered())

	var dst SomethingSet
	const txt = `[1, 2, 3, 3, 4, 5, 5, 6]`
	err := json.Unmarshal([]byte(txt), &dst)
	if assert.NoError(t, err) {
		assert.Equal(t, []Something{1, 2, 3, 4, 5, 6}, dst.Ordered())
	}

	buf, err := json.Marshal(set)
	if assert.NoError(t, err) {
		assert.Equal(t, `[1,2,3,4,5,6]`, string(buf))
	}
}
