package quiet

import (
	"context"
	"io"
	"sync"
	"time"
)

// Closer for when you need to coordinate delayed cleanup.
type Closer struct {
	closers []io.Closer
	logger  Logger
}

// enforce interface implementation
var _ io.Closer = (*Closer)(nil)

// SetLogger sets a custom logger for this closer.
// By default, all errors and panics are ignored.
func (c *Closer) SetLogger(logger Logger) {
	c.logger = logger
}

// Add multiple closers to be invoked by [Closer.CloseAll].
func (c *Closer) Add(closers ...io.Closer) {
	c.closers = append(c.closers, closers...)
}

// AddFunc adds a closer function to be invoked by [Closer.CloseAll].
func (c *Closer) AddFunc(closer func()) {
	c.closers = append(c.closers, &quietCloser{closer})
}

// AddFuncE adds a closer function to be invoked by [Closer.CloseAll].
func (c *Closer) AddFuncE(closer func() error) {
	c.closers = append(c.closers, &quietCloserE{closer})
}

// AddTimeout adds a closer function with a timeout to be invoked by [Closer.CloseAll].
func (c *Closer) AddTimeout(closer func(ctx context.Context) error, timeout time.Duration) {
	c.closers = append(c.closers, &timeoutCloser{close: closer, timeout: timeout})
}

// CloseAll will call each closer in reverse addition order.
// This mimics the invocation order of defers in a function.
func (c *Closer) CloseAll() {
	if c.logger == nil {
		c.logger = noopLogger{}
	}
	for i := len(c.closers) - 1; i >= 0; i-- {
		c.close(c.closers[i])
	}
}

// CloseAsync will attempt to call the closers in parallel.
// This is useful in speeding up an application's shutdown
// but does not guarantee an invocation order for closers.
func (c *Closer) CloseAsync() {
	if c.logger == nil {
		c.logger = noopLogger{}
	}
	var wg sync.WaitGroup
	wg.Add(len(c.closers))
	for i := len(c.closers) - 1; i >= 0; i-- {
		closer := c.closers[i]
		go func() {
			defer wg.Done()
			c.close(closer)
		}()
	}
	wg.Wait()
}

// Close for [io.Closer] compatibility.
// Calls [Closer.CloseAll] and returns nil.
func (c *Closer) Close() error {
	c.CloseAll()
	return nil
}

func (c *Closer) close(closer io.Closer) {
	defer func() {
		if p := recover(); p != nil {
			c.logger.Panic(p)
		}
	}()
	if err := closer.Close(); err != nil {
		c.logger.Error(err)
	}
}
