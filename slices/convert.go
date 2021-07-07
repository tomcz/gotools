package slices

import "github.com/cheekybits/genny/generic"

//go:generate genny -in=$GOFILE -out=convert_gen.go gen "SrcType=string,int,int64,uint64 DstType=interface{}"

type SrcType generic.Type
type DstType generic.Type

// ConvertSrcTypeToDstType generics helper.
func ConvertSrcTypeToDstType(src []SrcType) []DstType {
	dst := make([]DstType, len(src))
	for i, value := range src {
		dst[i] = value
	}
	return dst
}

// AppendSrcTypeToDstType generics helper.
func AppendSrcTypeToDstType(dst []DstType, src ...SrcType) []DstType {
	for _, value := range src {
		dst = append(dst, value)
	}
	return dst
}
