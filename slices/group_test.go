package slices

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	src := []string{
		"one_foo",
		"one_bar",
		"one_wee",
		"two_foo",
		"two_bar",
		"two_wee",
	}
	expected := map[string][]string{
		"one": {
			"one_foo",
			"one_bar",
			"one_wee",
		},
		"two": {
			"two_foo",
			"two_bar",
			"two_wee",
		},
	}
	groupFunc := func(value string) string {
		tokens := strings.Split(value, "_")
		return tokens[0]
	}
	assert.Equal(t, expected, GroupBy(src, groupFunc))
}
