package runner

import (
	"errors"
	"sync"
	"testing"

	"gotest.tools/v3/assert"
)

func TestRunner(t *testing.T) {
	var lock sync.Mutex
	var order []string

	app := New()
	app.Run(func() error {
		lock.Lock()
		order = append(order, "run")
		lock.Unlock()
		return errors.New("run error")
	})
	app.CleanupFunc(func() {
		lock.Lock()
		order = append(order, "cleanup")
		lock.Unlock()
	})
	err := app.Wait()

	assert.Error(t, err, "run error")
	assert.DeepEqual(t, order, []string{"run", "cleanup"})
}
