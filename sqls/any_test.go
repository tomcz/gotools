package sqls

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestMapToAny(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []any{1, 2, 3, 4, 5}
	assert.Equal(t, expected, MapToAny(input))
}
