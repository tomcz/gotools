package slices

import "github.com/cheekybits/genny/generic"

//go:generate genny -in=$GOFILE -out=split_gen.go gen "SliceType=string,int,int64,uint64"

type SliceType generic.Type

// SplitSliceType splits a slice into parts of a given length, with a remainder if necessary.
func SplitSliceType(src []SliceType, partLen int) [][]SliceType {
	srcLen := len(src)
	if srcLen == 0 {
		return nil
	}
	var dst [][]SliceType
	for a := 0; a < srcLen; a += partLen {
		z := a + partLen
		if z > srcLen {
			z = srcLen
		}
		dst = append(dst, src[a:z])
	}
	return dst
}
