package slices

import "github.com/cheekybits/genny/generic"

//go:generate genny -in=$GOFILE -out=convert_gen.go gen "Source=string,int,int64,uint64"

type Source generic.Type

// ConvertSourceToInterface generics helper.
func ConvertSourceToInterface(src []Source) []interface{} {
	dst := make([]interface{}, len(src))
	for i, value := range src {
		dst[i] = value
	}
	return dst
}

// AppendSourceToInterface generics helper.
func AppendSourceToInterface(dst []interface{}, src ...Source) []interface{} {
	for _, value := range src {
		dst = append(dst, value)
	}
	return dst
}
