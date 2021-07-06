package slices

import "github.com/cheekybits/genny/generic"

//go:generate genny -in=$GOFILE -out=convert_gen.go gen "Any=string,int,int64,uint64"

type Any generic.Type

// ConvertAnyToInterface generics helper.
func ConvertAnyToInterface(in []Any) []interface{} {
	out := make([]interface{}, len(in))
	for i, value := range in {
		out[i] = value
	}
	return out
}

// AppendAnyToInterface generics helper.
func AppendAnyToInterface(out []interface{}, in ...Any) []interface{} {
	for _, value := range in {
		out = append(out, value)
	}
	return out
}
