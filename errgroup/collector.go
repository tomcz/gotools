package errgroup

import "sync"

// Collect is a wrapper for the default behaviour of [Collector.Collect].
func Collect(funcs ...func() error) []error {
	return new(Collector).Collect(funcs...)
}

// Collector of errors.
type Collector struct {
	// Optional. The default handler will be used if this is nil.
	Handler PanicHandler
}

// Collect invokes all functions in panic-safe goroutines
// and returns any resulting errors or handled panics. The
// order of errors is undefined as it depends on the time
// of completion for each panic-safe goroutine.
func (c *Collector) Collect(funcs ...func() error) []error {
	if len(funcs) == 0 {
		return nil
	}
	if c.Handler == nil {
		c.Handler = defaultPanicHandler{}
	}
	errs := &errorList{}
	var wg sync.WaitGroup
	for _, f := range funcs {
		wg.Go(c.collect(errs, f))
	}
	wg.Wait()
	return errs.list
}

func (c *Collector) collect(errs *errorList, f func() error) func() {
	return func() {
		defer func() {
			if p := recover(); p != nil {
				errs.Append(c.Handler.Panic(p))
			}
		}()
		if err := f(); err != nil {
			errs.Append(err)
		}
	}
}

type errorList struct {
	list []error
	mux  sync.Mutex
}

func (e *errorList) Append(err error) {
	e.mux.Lock()
	e.list = append(e.list, err)
	e.mux.Unlock()
}
