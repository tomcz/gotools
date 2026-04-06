package errgroup

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// PanicGroup provides a panic handling wrapper around [golang.org/x/sync/errgroup.Group]
// to avoid application crashes when a goroutine encounters a panic.
type PanicGroup struct {
	group   *errgroup.Group
	handler PanicHandler
}

// New creates a [PanicGroup] without any context cancellation.
func New() *PanicGroup {
	return &PanicGroup{
		group:   &errgroup.Group{},
		handler: defaultPanicHandler{},
	}
}

// NewContext creates a [PanicGroup] with context cancellation.
// The returned context is cancelled on first error, first panic,
// or when the Wait function exits.
func NewContext(ctx context.Context) (*PanicGroup, context.Context) {
	group, ctx := errgroup.WithContext(ctx)
	pg := &PanicGroup{group: group, handler: defaultPanicHandler{}}
	return pg, ctx
}

// SetPanicHandler to provide a custom [PanicHandler] if you need to perform
// your own panic processing (e.g. send a notification to a tracking service).
func (g *PanicGroup) SetPanicHandler(handler PanicHandler) {
	if handler != nil {
		g.handler = handler
	} else {
		g.handler = defaultPanicHandler{}
	}
}

// Go delegates to [golang.org/x/sync/errgroup.Group.Go]
// and deals with any recovered panics using a [PanicHandler].
func (g *PanicGroup) Go(f func() error) {
	g.group.Go(func() (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = g.handler.Panic(p)
			}
		}()
		return f()
	})
}

// TryGo delegates to [golang.org/x/sync/errgroup.Group.TryGo]
// and deals with any recovered panics using a [PanicHandler].
func (g *PanicGroup) TryGo(f func() error) bool {
	return g.group.TryGo(func() (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = g.handler.Panic(p)
			}
		}()
		return f()
	})
}

// SetLimit delegates to [golang.org/x/sync/errgroup.Group.SetLimit].
func (g *PanicGroup) SetLimit(n int) {
	g.group.SetLimit(n)
}

// Wait delegates to [golang.org/x/sync/errgroup.Group.Wait].
func (g *PanicGroup) Wait() error {
	return g.group.Wait()
}
