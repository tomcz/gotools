package quiet

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	assert.True(t, closed, "close function was not called")
}

func TestClose_Error(t *testing.T) {
	defer SetLogger(nil)

	log := &testLogger{}
	SetLogger(log)

	testErr := errors.New("test error")
	CloseFuncE(func() error { return testErr })
	assert.Equal(t, testErr, log.err)
}

func TestClose_Panic(t *testing.T) {
	defer SetLogger(nil)

	log := &testLogger{}
	SetLogger(log)

	testErr := errors.New("test error")
	CloseFuncE(func() error { panic(testErr) })
	assert.Equal(t, testErr, log.panic)
}

func TestClose_Timeout(t *testing.T) {
	closed := false
	close := func(context.Context) error {
		closed = true
		return nil
	}
	CloseWithTimeout(close, time.Minute)
	assert.True(t, closed, "close function was not called")
}
