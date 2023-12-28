package maths

import "cmp"

// Min returns the least value.
func Min[X cmp.Ordered](a, b X) X {
	if a < b {
		return a
	}
	return b
}

// MinOf returns the least value.
func MinOf[X cmp.Ordered](a X, bb ...X) X {
	for _, b := range bb {
		a = Min(a, b)
	}
	return a
}
