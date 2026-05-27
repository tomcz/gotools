package errgroup

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestCollect(t *testing.T) {
	collector := &Collector{Handler: testPanicHandler{}}
	res := collector.Collect(
		func() error { return nil },
		func() error { return errors.New("wibble") },
		func() error { panic("d'oh") },
	)
	assert.Equal(t, len(res), 2)
	assert.Assert(t, slices.ContainsFunc(res, func(err error) bool {
		return strings.Contains(err.Error(), "wibble")
	}))
	assert.Assert(t, slices.ContainsFunc(res, func(err error) bool {
		return strings.Contains(err.Error(), "d'oh crap")
	}))
}

func TestCollect_Default(t *testing.T) {
	collector := &Collector{}
	res := collector.Collect(
		func() error { return nil },
		func() error { return errors.New("wibble") },
		func() error { panic("d'oh") },
	)
	assert.Equal(t, len(res), 2)
	assert.Assert(t, slices.ContainsFunc(res, func(err error) bool {
		return strings.Contains(err.Error(), "wibble")
	}))
	assert.Assert(t, slices.ContainsFunc(res, func(err error) bool {
		return strings.Contains(err.Error(), "panic: d'oh")
	}))
}
