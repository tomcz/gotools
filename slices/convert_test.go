package slices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	expected := []any{1, 2, 3, 4}
	actual := ConvertSourceToAny([]Source{1, 2, 3, 4})
	assert.Equal(t, expected, actual)
}

func TestAppend(t *testing.T) {
	expected := []any{"foo", 2, 3, 4}
	actual := AppendSourceToAny([]any{"foo"}, 2, 3, 4)
	assert.Equal(t, expected, actual)
}
