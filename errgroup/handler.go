package errgroup

import (
	"fmt"
	"runtime/debug"
)

// PanicHandler to deal with recovered panics.
type PanicHandler interface {
	Panic(p any) error
}

// enforce interface implementation
var _ PanicHandler = (*DefaultPanicHandler)(nil)

// DefaultPanicHandler used by both [PanicGroup] and [Collector].
type DefaultPanicHandler struct{}

// Panic combines the recovered panic and the current [debug.Stack] into a single error.
func (d DefaultPanicHandler) Panic(p any) error {
	stack := string(debug.Stack())
	if err, ok := p.(error); ok {
		return fmt.Errorf("panic: %w; stack: %s", err, stack)
	}
	return fmt.Errorf("panic: %+v; stack: %s", p, stack)
}
