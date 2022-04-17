package errgroup

import (
	"context"
	"fmt"
	"runtime/debug"

	"golang.org/x/sync/errgroup"
)

// Group provides an interface compatible with golang.org/x/sync/errgroup
// for instances that enhance the capabilities of Groups.
type Group interface {
	Go(f func() error)
	Wait() error
}

type panicGroup struct {
	group *errgroup.Group
}

// New creates a panic-handling Group,
// without any context cancellation.
func New() Group {
	group := &errgroup.Group{}
	return &panicGroup{group}
}

// WithContext creates a panic-handling Group.
// The returned context is cancelled on first error,
// first panic, or when the Wait function exits.
func WithContext(ctx context.Context) (Group, context.Context) {
	group, ctx := errgroup.WithContext(ctx)
	return &panicGroup{group}, ctx
}

func (p *panicGroup) Wait() error {
	return p.group.Wait()
}

func (p *panicGroup) Go(f func() error) {
	p.group.Go(func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				if ex, ok := r.(error); ok {
					err = fmt.Errorf("panic: %w; stack: %s", ex, stack)
				} else {
					err = fmt.Errorf("panic: %v; stack: %s", r, stack)
				}
			}
		}()
		return f()
	})
}
