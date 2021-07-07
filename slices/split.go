package slices

import "github.com/cheekybits/genny/generic"

//go:generate genny -in=$GOFILE -out=split_gen.go gen "Slice=string,int,int64,uint64"

type Slice generic.Type

// SplitSlice splits a slice into parts of a given length, with a remainder if necessary.
func SplitSlice(src []Slice, partLen int) [][]Slice {
	srcLen := len(src)
	if srcLen == 0 {
		return nil
	}
	var dst [][]Slice
	for a := 0; a < srcLen; a += partLen {
		z := a + partLen
		if z > srcLen {
			z = srcLen
		}
		dst = append(dst, src[a:z])
	}
	return dst
}
