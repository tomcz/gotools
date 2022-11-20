package sqls

import "github.com/tomcz/gotools/slices"

// MapToAny performs a type conversion useful for query arguments.
func MapToAny[A any](src []A) []any {
	return slices.Map(src, func(a A) any { return a })
}
