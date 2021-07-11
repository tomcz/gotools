// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package slices

// ConvertUint64ToInterface generics helper.
func ConvertUint64ToInterface(src []uint64) []interface{} {
	dst := make([]interface{}, len(src))
	for i, value := range src {
		dst[i] = value
	}
	return dst
}

// AppendUint64ToInterface generics helper.
func AppendUint64ToInterface(dst []interface{}, src ...uint64) []interface{} {
	for _, value := range src {
		dst = append(dst, value)
	}
	return dst
}
