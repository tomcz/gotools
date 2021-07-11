// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package slices

// ConvertStringToInterface generics helper.
func ConvertStringToInterface(src []string) []interface{} {
	dst := make([]interface{}, len(src))
	for i, value := range src {
		dst[i] = value
	}
	return dst
}

// AppendStringToInterface generics helper.
func AppendStringToInterface(dst []interface{}, src ...string) []interface{} {
	for _, value := range src {
		dst = append(dst, value)
	}
	return dst
}
