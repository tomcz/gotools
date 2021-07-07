package slices

import "github.com/cheekybits/genny/generic"

//go:generate genny -in=$GOFILE -out=split_gen.go gen "Value=string,int,int64,uint64"

type Value generic.Type

// SplitValue splits a slice into parts of a given length, with a remainder if necessary.
func SplitValue(src []Value, partLen int) [][]Value {
	srcLen := len(src)
	if srcLen == 0 {
		return nil
	}
	var dst [][]Value
	for a := 0; a < srcLen; a += partLen {
		z := a + partLen
		if z > srcLen {
			z = srcLen
		}
		dst = append(dst, src[a:z])
	}
	return dst
}
