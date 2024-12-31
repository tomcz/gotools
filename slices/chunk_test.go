package slices

import (
	"testing"
	"testing/quick"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestChunk(t *testing.T) {
	tests := []struct {
		name     string
		numParts int
		src      []int
		expected [][]int
	}{
		{
			name:     "empty slice",
			numParts: 2,
			src:      nil,
			expected: nil,
		},
		{
			name:     "zero parts",
			numParts: 0,
			src:      []int{1, 2, 3},
			expected: nil,
		},
		{
			name:     "single part",
			numParts: 1,
			src:      []int{1, 2, 3},
			expected: [][]int{{1, 2, 3}},
		},
		{
			name:     "srcLen equals numParts",
			numParts: 3,
			src:      []int{1, 2, 3},
			expected: [][]int{{1}, {2}, {3}},
		},
		{
			name:     "srcLen is less than numParts",
			numParts: 4,
			src:      []int{1, 2, 3},
			expected: [][]int{{1}, {2}, {3}, nil},
		},
		{
			name:     "even chunks",
			numParts: 2,
			src: []int{
				1, 2, 3,
				4, 5, 6,
				7, 8, 9,
				10,
			},
			expected: [][]int{
				{1, 2, 3, 4, 5},
				{6, 7, 8, 9, 10},
			},
		},
		{
			name:     "uneven chunks",
			numParts: 3,
			src: []int{
				1, 2, 3,
				4, 5, 6,
				7, 8, 9,
				10,
			},
			expected: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9, 10},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.DeepEqual(t, test.expected, Chunk(test.src, test.numParts))
		})
	}
}

func TestChunk_fuzz(t *testing.T) {
	fn := func(src []int, numParts uint8) bool {
		chunks := Chunk(src, int(numParts))
		if len(chunks) == 0 && (len(src) == 0 || numParts == 0) {
			return true
		}
		if !assert.Check(t, is.Len(chunks, int(numParts))) {
			return false
		}
		var joined []int
		var minLen, maxLen int
		for i, chunk := range chunks {
			joined = append(joined, chunk...)
			chunkLen := len(chunk)
			if i == 0 {
				minLen = chunkLen
				maxLen = chunkLen
			} else {
				minLen = min(chunkLen, minLen)
				maxLen = max(chunkLen, maxLen)
			}
		}
		if !assert.Check(t, maxLen-minLen < 2, "maxLen %d, minLen %d", maxLen, minLen) {
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
