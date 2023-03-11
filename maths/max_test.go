package maths

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestMax(t *testing.T) {
	assert.Equal(t, 2, Max(1, 2))
	assert.Equal(t, 2, Max(2, 1))
	assert.Equal(t, 2, Max(2, 2))
}

func TestMaxOf(t *testing.T) {
	assert.Equal(t, 3, MaxOf(1, 2, 3))
	assert.Equal(t, 3, MaxOf(3, 2, 1))
	assert.Equal(t, 3, MaxOf(3, 3, 3))
	assert.Equal(t, 3, MaxOf(3))
}
