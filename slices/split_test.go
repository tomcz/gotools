package slices

import (
	"testing"
	"testing/quick"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
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

func TestSplit_fuzz(t *testing.T) {
	fn := func(src []int, partLen uint8) bool {
		splits := Split(src, int(partLen))
		if len(splits) == 0 && (len(src) == 0 || partLen == 0) {
			return true
		}
		var joined []int
		var minLen, maxLen int
		for i, split := range splits {
			joined = append(joined, split...)
			splitLen := len(split)
			if i == 0 {
				minLen = splitLen
				maxLen = splitLen
			} else {
				minLen = min(splitLen, minLen)
				maxLen = max(splitLen, maxLen)
			}
		}
		if !assert.Check(t, minLen <= int(partLen), "minLen %d, partLen %d", minLen, partLen) {
			return false
		}
		if !assert.Check(t, maxLen <= int(partLen), "maxLen %d, partLen %d", maxLen, partLen) {
			return false
		}
		if !assert.Check(t, is.DeepEqual(src, joined)) {
			return false
		}
		return true
	}
	if err := quick.Check(fn, nil); err != nil {
		t.Error(err)
	}
}
