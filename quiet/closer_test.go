package quiet

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
	"time"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestCloserOrder(t *testing.T) {
	var called []string
	c1 := func() {
		called = append(called, "1")
	}
	c2 := func() {
		called = append(called, "2")
	}
	c3 := func() {
		called = append(called, "3")
	}
	c4 := func() error {
		called = append(called, "4")
		return nil
	}
	c5 := func(context.Context) error {
		called = append(called, "5")
		return nil
	}

	c1w := &quietCloser{close: c1}
	c2w := &quietCloser{close: c2}

	closer := &Closer{}
	closer.Add(c1w, c2w)
	closer.AddFunc(c3)
	closer.AddFuncE(c4)
	closer.AddTimeout(c5, time.Minute)

	assert.NilError(t, closer.Close())
	assert.DeepEqual(t, []string{"5", "4", "3", "2", "1"}, called)
}

func TestCloseAsync(t *testing.T) {
	var mux sync.Mutex
	var called []string
	c1 := func() {
		mux.Lock()
		called = append(called, "1")
		mux.Unlock()
	}
	c2 := func() {
		mux.Lock()
		called = append(called, "2")
		mux.Unlock()
	}
	c3 := func() {
		mux.Lock()
		called = append(called, "3")
		mux.Unlock()
	}
	c4 := func() error {
		mux.Lock()
		called = append(called, "4")
		mux.Unlock()
		return nil
	}
	c5 := func(context.Context) error {
		mux.Lock()
		called = append(called, "5")
		mux.Unlock()
		return nil
	}

	c1w := &quietCloser{close: c1}
	c2w := &quietCloser{close: c2}

	closer := &Closer{}
	closer.Add(c1w, c2w)
	closer.AddFunc(c3)
	closer.AddFuncE(c4)
	closer.AddTimeout(c5, time.Minute)

	closer.CloseAsync()
	slices.Sort(called)

	assert.DeepEqual(t, []string{"1", "2", "3", "4", "5"}, called)
}

func TestCloserLogger(t *testing.T) {
	logger := &Collector{}

	closer := &Closer{}
	closer.SetLogger(logger)

	testPanic := errors.New("test panic")
	closer.AddFunc(func() { panic(testPanic) })

	testError := errors.New("test error")
	closer.AddFuncE(func() error { return testError })

	closer.CloseAll()

	assert.Assert(t, !logger.IsEmpty())
	assert.Assert(t, is.Contains(logger.Panics, testPanic))
	assert.Assert(t, is.Contains(logger.Errors, testError))
}
