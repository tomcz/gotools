package quiet

import "io"

type noopCloser struct{}

// Noop returns a closer that does nothing.
func Noop() io.Closer {
	return noopCloser{}
}

func (n noopCloser) Close() error {
	return nil
}
