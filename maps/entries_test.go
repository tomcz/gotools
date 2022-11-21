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

	expected := []Entry[string, int]{
		{Key: "1", Val: 1},
		{Key: "2", Val: 2},
		{Key: "3", Val: 3},
		{Key: "4", Val: 4},
		{Key: "5", Val: 5},
	}
	assert.Equal(t, expected, SortedEntries(data))
	assert.Equal(t, data, FromEntries(expected))
}
