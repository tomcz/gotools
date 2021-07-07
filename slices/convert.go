package slices

import "github.com/cheekybits/genny/generic"

//go:generate genny -in=$GOFILE -out=convert_gen.go gen "Src=string,int,int64,uint64 Dst=interface{}"

type Src generic.Type
type Dst generic.Type

// ConvertSrcToDst generics helper.
func ConvertSrcToDst(src []Src) []Dst {
	dst := make([]Dst, len(src))
	for i, value := range src {
		dst[i] = value
	}
	return dst
}

// AppendSrcToDst generics helper.
func AppendSrcToDst(dst []Dst, src ...Src) []Dst {
	for _, value := range src {
		dst = append(dst, value)
	}
	return dst
}
