package maths

import "golang.org/x/exp/constraints"

// Max returns the greatest value.
func Max[X constraints.Ordered](a, b X) X {
	if a > b {
		return a
	}
	return b
}

// MaxOf returns the greatest value.
func MaxOf[X constraints.Ordered](a X, bb ...X) X {
	for _, b := range bb {
		a = Max(a, b)
	}
	return a
}
