package maths

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestMin(t *testing.T) {
	assert.Equal(t, 1, Min(1, 2))
	assert.Equal(t, 1, Min(2, 1))
	assert.Equal(t, 1, Min(1, 1))
}

func TestMinOf(t *testing.T) {
	assert.Equal(t, 1, MinOf(1, 2, 3))
	assert.Equal(t, 1, MinOf(3, 2, 1))
	assert.Equal(t, 1, MinOf(1, 1, 1))
	assert.Equal(t, 1, MinOf(1))
}
