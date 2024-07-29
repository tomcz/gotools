package quiet

import (
	"context"
	"errors"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

type testLogger struct {
	err   error
	panic any
}

func (l *testLogger) Error(err error) {
	l.err = err
}

func (l *testLogger) Panic(p any) {
	l.panic = p
}

func TestClose_Quiet(t *testing.T) {
	closed := false
	CloseFunc(func() { closed = true })
	assert.Assert(t, closed, "close function was not called")
}

func TestClose_Error(t *testing.T) {
	defer SetLogger(nil)

	logger := &testLogger{}
	SetLogger(logger)

	testErr := errors.New("test error")
	CloseFuncE(func() error { return testErr })
	assert.Equal(t, testErr, logger.err)
}

func TestClose_Panic(t *testing.T) {
	defer SetLogger(nil)

	logger := &testLogger{}
	SetLogger(logger)

	testErr := errors.New("test error")
	CloseFuncE(func() error { panic(testErr) })
	assert.Equal(t, testErr, logger.panic)
}

func TestClose_Timeout(t *testing.T) {
	closed := false
	closer := func(context.Context) error {
		closed = true
		return nil
	}
	CloseWithTimeout(closer, time.Minute)
	assert.Assert(t, closed, "close function was not called")
}
