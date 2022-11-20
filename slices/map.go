package slices

// Map converts from one type to another.
func Map[A any, B any](src []A, mapper func(A) B) []B {
	dest := make([]B, len(src))
	for i, val := range src {
		dest[i] = mapper(val)
	}
	return dest
}

// MapErr allows the conversion to fail for any value that returns an error.
func MapErr[A any, B any](src []A, mapper func(A) (B, error)) ([]B, error) {
	dest := make([]B, len(src))
	for i, in := range src {
		out, err := mapper(in)
		if err != nil {
			return nil, err
		}
		dest[i] = out
	}
	return dest, nil
}
