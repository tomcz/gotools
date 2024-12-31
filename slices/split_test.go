package slices

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		name     string
		src      []int
		partLen  int
		expected [][]int
	}{
		{
			name:     "empty slice",
			partLen:  10,
			src:      nil,
			expected: nil,
		},
		{
			name:     "partLen is zero",
			partLen:  0,
			src:      []int{1, 2, 3, 4, 5},
			expected: nil,
		},
		{
			name:     "partLen is one",
			partLen:  1,
			src:      []int{1, 2, 3, 4, 5},
			expected: [][]int{{1}, {2}, {3}, {4}, {5}},
		},
		{
			name:     "srcLen is less than partLen",
			partLen:  10,
			src:      []int{1, 2, 3, 4, 5},
			expected: [][]int{{1, 2, 3, 4, 5}},
		},
		{
			name:     "srcLen equals partLen",
			partLen:  5,
			src:      []int{1, 2, 3, 4, 5},
			expected: [][]int{{1, 2, 3, 4, 5}},
		},
		{
			name:    "even split",
			partLen: 3,
			src: []int{
				1, 2, 3,
				4, 5, 6,
				7, 8, 9,
			},
			expected: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
		},
		{
			name:    "uneven split",
			partLen: 3,
			src: []int{
				1, 2, 3,
				4, 5, 6,
				7, 8, 9,
				10,
			},
			expected: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
				{10},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.DeepEqual(t, test.expected, Split(test.src, test.partLen))
		})
	}
}
