package sqls

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestMapToAny(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []any{1, 2, 3, 4, 5}
	assert.DeepEqual(t, expected, MapToAny(input))
}
