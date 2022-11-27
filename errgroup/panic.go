package errgroup

import (
	"context"
	"fmt"
	"runtime/debug"

	"golang.org/x/sync/errgroup"
)

// PanicHandler processes the recovered panic.
type PanicHandler func(p any) error

func defaultPanicHandler(p any) error {
	stack := string(debug.Stack())
	if ex, ok := p.(error); ok {
		return fmt.Errorf("panic: %w; stack: %s", ex, stack)
	}
	return fmt.Errorf("panic: %+v; stack: %s", p, stack)
}

// Group provides an interface compatible with golang.org/x/sync/errgroup
// for instances that enhance the capabilities of Groups.
type Group interface {
	Go(f func() error)
	Wait() error
	OnPanic(ph PanicHandler)
}

type panicGroup struct {
	group *errgroup.Group
	panic PanicHandler
}

// New creates a panic-handling Group,
// without any context cancellation.
func New() Group {
	group := &errgroup.Group{}
	return &panicGroup{
		group: group,
		panic: defaultPanicHandler,
	}
}

// WithContext creates a panic-handling Group.
// The returned context is cancelled on first error,
// first panic, or when the Wait function exits.
func WithContext(ctx context.Context) (Group, context.Context) {
	group, ctx := errgroup.WithContext(ctx)
	pg := &panicGroup{
		group: group,
		panic: defaultPanicHandler,
	}
	return pg, ctx
}

func (p *panicGroup) OnPanic(ph PanicHandler) {
	if ph != nil {
		p.panic = ph
	}
}

func (p *panicGroup) Wait() error {
	return p.group.Wait()
}

func (p *panicGroup) Go(f func() error) {
	p.group.Go(func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = p.panic(r)
			}
		}()
		return f()
	})
}
