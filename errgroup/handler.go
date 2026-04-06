package errgroup

import (
	"fmt"
	"runtime/debug"
)

// PanicHandler to deal with recovered panics. The default handler combines
// a recovered panic and the current [debug.Stack] into a single error.
type PanicHandler interface {
	Panic(p any) error
}

var _ PanicHandler = (*defaultPanicHandler)(nil)

type defaultPanicHandler struct{}

func (d defaultPanicHandler) Panic(p any) error {
	stack := string(debug.Stack())
	if err, ok := p.(error); ok {
		return fmt.Errorf("panic: %w; stack: %s", err, stack)
	}
	return fmt.Errorf("panic: %+v; stack: %s", p, stack)
}
