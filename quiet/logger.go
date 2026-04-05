package quiet

// Logger allows logging of errors and panics.
type Logger interface {
	Error(err error)
	Panic(p any)
}

type noopLogger struct{}

func (n noopLogger) Error(error) {}

func (n noopLogger) Panic(any) {}

// Collector is a [Logger] variant that collects all
// encountered errors and panics for later review.
type Collector struct {
	Errors []error
	Panics []any
}

// enforce interface implementation
var _ Logger = (*Collector)(nil)

// Error appends the error to the [Collector] Errors slice.
func (c *Collector) Error(err error) {
	c.Errors = append(c.Errors, err)
}

// Panic appends the panic to the [Collector] Panics slice.
func (c *Collector) Panic(p any) {
	c.Panics = append(c.Panics, p)
}

// IsEmpty returns true if there aren't any errors or panics.
func (c *Collector) IsEmpty() bool {
	return len(c.Errors)+len(c.Panics) == 0
}
