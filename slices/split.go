package slices

// Split given slice into equal sized parts.
// The last part may have an unequal length
// if the size of the source slice is not
// divisible by the required length without
// leaving a remainder.
func Split[X any](src []X, partLen int) [][]X {
	srcLen := len(src)
	if srcLen == 0 {
		return nil
	}
	var dst [][]X
	for a := 0; a < srcLen; a += partLen {
		z := a + partLen
		if z > srcLen {
			z = srcLen
		}
		dst = append(dst, src[a:z])
	}
	return dst
}
