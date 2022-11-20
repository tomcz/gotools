package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortedEntries(t *testing.T) {
	data := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
	}
	assert.Equal(t, []string{"1", "2", "3", "4", "5"}, SortedKeys(data))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, SortedValues(data))
}
