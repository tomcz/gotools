package split

// Interface splits a slice into parts of a given length, with a remainder if necessary.
func Interface(in []interface{}, partLen int) [][]interface{} {
	inLen := len(in)
	if inLen == 0 {
		return nil
	}
	var out [][]interface{}
	for a := 0; a < inLen; a += partLen {
		z := a + partLen
		if z > inLen {
			z = inLen
		}
		out = append(out, in[a:z])
	}
	return out
}
