package quiet

import (
	"errors"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestLogger(t *testing.T) {
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
