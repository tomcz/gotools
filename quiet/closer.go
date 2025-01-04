package quiet

import (
	"context"
	"io"
	"time"
)

// Closer for when you need to coordinate delayed cleanup.
type Closer struct {
	closers []io.Closer
}

// enforce interface implementation
var _ io.Closer = &Closer{}

// Add multiple closers to be invoked by CloseAll.
//
// See also Close.
func (c *Closer) Add(closers ...io.Closer) {
	c.closers = append(c.closers, closers...)
}

// AddFunc adds a closer function to be invoked by CloseAll.
//
// See also CloseFunc.
func (c *Closer) AddFunc(closer func()) {
	c.closers = append(c.closers, &quietCloser{closer})
}

// AddFuncE adds a closer function to be invoked by CloseAll.
//
// See also CloseFuncE.
func (c *Closer) AddFuncE(closer func() error) {
	c.closers = append(c.closers, &quietCloserE{closer})
}

// AddTimeout adds a closer function with a timeout to be invoked by CloseAll.
//
// See also CloseWithTimeout.
func (c *Closer) AddTimeout(closer func(ctx context.Context) error, timeout time.Duration) {
	c.closers = append(c.closers, &timeoutCloser{close: closer, timeout: timeout})
}

// CloseAll will call each closer in reverse addition order.
// This mimics the invocation order of defers in a function.
func (c *Closer) CloseAll() {
	for i := len(c.closers) - 1; i >= 0; i-- {
		Close(c.closers[i])
	}
}

// Close for io.Closer compatibility.
// Calls CloseAll and returns nil.
func (c *Closer) Close() error {
	c.CloseAll()
	return nil
}
