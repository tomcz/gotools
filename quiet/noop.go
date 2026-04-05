package quiet

import "io"

type noopCloser struct{}

// Noop returns an [io.Closer] that does nothing.
func Noop() io.Closer {
	return noopCloser{}
}

func (n noopCloser) Close() error {
	return nil
}
