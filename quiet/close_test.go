package quiet

import (
	"errors"
	"testing"

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
