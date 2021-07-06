package chunks

import "github.com/cheekybits/genny/generic"

//go:generate genny -in=$GOFILE -out=split.go gen "Something=string,int,int64,uint64"

type Something generic.Type

// SplitSomething splits a slice into parts of a given length, with a remainder if necessary.
func SplitSomething(in []Something, partLen int) [][]Something {
	inLen := len(in)
	if inLen == 0 {
		return nil
	}
	var out [][]Something
	for a := 0; a < inLen; a += partLen {
		z := a + partLen
		if z > inLen {
			z = inLen
		}
		out = append(out, in[a:z])
	}
	return out
}
