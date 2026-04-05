package errgroup

import (
	"context"
	"fmt"
	"runtime/debug"

	"golang.org/x/sync/errgroup"
)

// Group provides an interface compatible with [golang.org/x/sync/errgroup.Group].
type Group interface {
	Go(f func() error)
	TryGo(f func() error) bool
	SetLimit(n int)
	Wait() error
}

var _ Group = (*errgroup.Group)(nil)
var _ Group = (*panicGroup)(nil)

// PanicHandler processes the recovered panic.
type PanicHandler func(p any) error

// Opt is a configuration option.
type Opt func(g *panicGroup)

type panicGroup struct {
	group  *errgroup.Group
	handle PanicHandler
}

// WithPanicHandler overrides the default panic handler.
func WithPanicHandler(ph PanicHandler) Opt {
	return func(p *panicGroup) {
		p.handle = ph
	}
}

// New creates a panic-handling [Group] without any context cancellation.
func New(opts ...Opt) Group {
	group := &errgroup.Group{}
	pg := &panicGroup{group: group}
	pg.configure(opts)
	return pg
}

// NewContext creates a panic-handling [Group].
// The returned context is cancelled on first error,
// first panic, or when the Wait function exits.
func NewContext(ctx context.Context, opts ...Opt) (Group, context.Context) {
	group, ctx := errgroup.WithContext(ctx)
	pg := &panicGroup{group: group}
	pg.configure(opts)
	return pg, ctx
}

func (pg *panicGroup) configure(opts []Opt) {
	for _, opt := range opts {
		opt(pg)
	}
	if pg.handle == nil {
		pg.handle = defaultPanicHandler
	}
}

func (pg *panicGroup) Go(f func() error) {
	pg.group.Go(func() (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = pg.handle(p)
			}
		}()
		return f()
	})
}

func (pg *panicGroup) TryGo(f func() error) bool {
	return pg.group.TryGo(func() (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = pg.handle(p)
			}
		}()
		return f()
	})
}

func (pg *panicGroup) SetLimit(n int) {
	pg.group.SetLimit(n)
}

func (pg *panicGroup) Wait() error {
	return pg.group.Wait()
}

func defaultPanicHandler(p any) error {
	stack := string(debug.Stack())
	if err, ok := p.(error); ok {
		return fmt.Errorf("panic: %w\nstack: %s", err, stack)
	}
	return fmt.Errorf("panic: %+v\nstack: %s", p, stack)
}
