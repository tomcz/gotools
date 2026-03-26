package quiet

// Logger allows logging of errors and panics.
type Logger interface {
	Error(err error)
	Panic(p any)
}

var log Logger = noopLogger{}

// SetLogger sets a panic & error logger for the package
// rather than the default noop logger. Passing in a nil
// logger will reset the package logger to default.
func SetLogger(logger Logger) {
	if logger == nil {
		log = noopLogger{}
	} else {
		log = logger
	}
}

type noopLogger struct{}

func (n noopLogger) Error(error) {}

func (n noopLogger) Panic(any) {}
