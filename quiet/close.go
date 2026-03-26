package quiet

import (
	"context"
	"io"
	"time"
)

// Close quietly invokes the closer.
// Any errors or panics will be ignored.
func Close(closer io.Closer) {
	defer func() {
		_ = recover()
	}()
	_ = closer.Close()
}

// CloseFunc quietly invokes the given function.
// Any errors or panics will be ignored.
func CloseFunc(close func()) {
	Close(&quietCloser{close})
}

// CloseFuncE quietly invokes the given function.
// Any errors or panics will be ignored.
func CloseFuncE(close func() error) {
	Close(&quietCloserE{close})
}

// CloseWithTimeout quietly invokes the given function with the timeout set on its context.
// Any errors or panics will be ignored.
func CloseWithTimeout(close func(ctx context.Context) error, timeout time.Duration) {
	Close(&timeoutCloser{close: close, timeout: timeout})
}

type quietCloserE struct {
	close func() error
}

func (c *quietCloserE) Close() error {
	return c.close()
}

type quietCloser struct {
	close func()
}

func (c *quietCloser) Close() error {
	c.close()
	return nil
}

type timeoutCloser struct {
	close   func(context.Context) error
	timeout time.Duration
}

func (c *timeoutCloser) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.close(ctx)
}
