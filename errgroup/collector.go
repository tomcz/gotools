package errgroup

import "sync"

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

	var handler PanicHandler
	if c.Handler != nil {
		handler = c.Handler
	} else {
		handler = defaultPanicHandler{}
	}

	res := &errorList{}
	var wg sync.WaitGroup
	wg.Add(len(funcs))

	run := func(f func() error) {
		defer func() {
			if p := recover(); p != nil {
				res.Append(handler.Panic(p))
			}
			wg.Done()
		}()
		if err := f(); err != nil {
			res.Append(err)
		}
	}
	for _, f := range funcs {
		go run(f)
	}

	wg.Wait()
	return res.list
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
