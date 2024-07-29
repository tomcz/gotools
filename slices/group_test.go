package slices

import (
	"errors"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
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
	assert.DeepEqual(t, expected, GroupBy(src, groupFunc))

	groupFuncOk := func(value string) (string, error) {
		return groupFunc(value), nil
	}
	actual, err := GroupByErr(src, groupFuncOk)
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)

	groupFuncBad := func(value string) (string, error) {
		return "", errors.New("test error")
	}
	actual, err = GroupByErr(src, groupFuncBad)
	assert.Error(t, err, "test error")
	assert.Assert(t, is.Nil(actual))
}
