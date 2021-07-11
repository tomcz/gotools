package slices

import "github.com/cheekybits/genny/generic"

//go:generate genny -in=$GOFILE -out=convert_int.go     gen "Source=int"
//go:generate genny -in=$GOFILE -out=convert_int64.go   gen "Source=int64"
//go:generate genny -in=$GOFILE -out=convert_uint64.go  gen "Source=uint64"
//go:generate genny -in=$GOFILE -out=convert_strings.go gen "Source=string"

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
