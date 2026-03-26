package quiet

import (
	"context"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestClose_Func(t *testing.T) {
	closed := false
	CloseFunc(func() { closed = true })
	assert.Assert(t, closed, "close function was not called")
}

func TestClose_FuncE(t *testing.T) {
	closed := false
	CloseFuncE(func() error { closed = true; return nil })
	assert.Assert(t, closed, "close function was not called")
}

func TestClose_Timeout(t *testing.T) {
	closed := false
	closer := func(context.Context) error { closed = true; return nil }
	CloseWithTimeout(closer, time.Minute)
	assert.Assert(t, closed, "close function was not called")
}

func TestClose_Panic(t *testing.T) {
	CloseFunc(func() { panic("test panic") })
}
