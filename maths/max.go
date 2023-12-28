package maths

import "cmp"

// Max returns the greatest value.
func Max[X cmp.Ordered](a, b X) X {
	if a > b {
		return a
	}
	return b
}

// MaxOf returns the greatest value.
func MaxOf[X cmp.Ordered](a X, bb ...X) X {
	for _, b := range bb {
		a = Max(a, b)
	}
	return a
}
