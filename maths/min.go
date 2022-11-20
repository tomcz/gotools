package maths

import "golang.org/x/exp/constraints"

// Min returns the least value.
func Min[X constraints.Ordered](a, b X) X {
	if a < b {
		return a
	}
	return b
}

// MinOf returns the least value.
func MinOf[X constraints.Ordered](a X, bb ...X) X {
	for _, b := range bb {
		a = Min(a, b)
	}
	return a
}
