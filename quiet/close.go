package quiet

import (
	"context"
	"io"
	"time"
)

// Close quietly invokes the closer.
// Any errors or panics will be logged by this package's Logger.
func Close(closer io.Closer) {
	defer func() {
		if p := recover(); p != nil {
			log.Panic(p)
		}
	}()
	if err := closer.Close(); err != nil {
		log.Error(err)
	}
}

// CloseFunc quietly invokes the given function.
// Any panics will be logged by this package's Logger.
func CloseFunc(close func()) {
	Close(&quietCloser{close})
}

// CloseFuncE quietly invokes the given function.
// Any errors or panics will be logged by this package's Logger.
func CloseFuncE(close func() error) {
	Close(&quietCloserE{close})
}

// CloseWithTimeout quiety invokes the given function with the timeout set on its context.
// Any errors or panics will be logged by this package's Logger.
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
