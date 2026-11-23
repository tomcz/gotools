package runner

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"gotest.tools/v3/assert"
)

func TestRunner(t *testing.T) {
	var lock sync.Mutex
	var order []string

	app := New(WithAsyncShutdown())
	app.Run(func() error {
		lock.Lock()
		order = append(order, "run")
		lock.Unlock()
		return nil
	})
	app.ShutdownFunc(func() {
		lock.Lock()
		order = append(order, "shutdown")
		lock.Unlock()
	})
	err := app.Wait()

	assert.NilError(t, err)
	assert.DeepEqual(t, order, []string{"run", "shutdown"})
}

func TestRunner_RunError(t *testing.T) {
	var cleanedUp atomic.Bool

	app := New()
	app.Run(func() error {
		return errors.New("test error")
	})
	app.ShutdownFunc(func() {
		cleanedUp.Store(true)
	})
	err := app.Wait()

	assert.Error(t, err, "test error")
	assert.Assert(t, cleanedUp.Load())
}

func TestRunner_RunPanic(t *testing.T) {
	var cleanedUp atomic.Bool

	app := New()
	app.Run(func() error {
		panic("test panic")
	})
	app.ShutdownFunc(func() {
		cleanedUp.Store(true)
	})
	err := app.Wait()

	assert.ErrorContains(t, err, "test panic")
	assert.Assert(t, cleanedUp.Load())
}
