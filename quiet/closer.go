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

// Add multiple closers to be invoked by Close.
func (c *Closer) Add(closers ...io.Closer) {
	c.closers = append(c.closers, closers...)
}

// AddFunc adds a cleanup function to be invoked by Close.
func (c *Closer) AddFunc(close func()) {
	c.closers = append(c.closers, &quietCloser{close})
}

// AddFuncE adds a cleanup function to be invoked by Close.
func (c *Closer) AddFuncE(close func() error) {
	c.closers = append(c.closers, &quietCloserE{close})
}

// AddTimeout adds a shutdown function to be invoked by Close.
func (c *Closer) AddTimeout(close func(ctx context.Context) error, timeout time.Duration) {
	c.closers = append(c.closers, &timeoutCloser{close: close, timeout: timeout})
}

// Close will call each closer in reverse order to addition.
// This mimics the invocation order of defers in a function.
func (c *Closer) Close() error {
	for i := len(c.closers) - 1; i >= 0; i-- {
		Close(c.closers[i])
	}
	return nil
}