package slices

// Chunk given slice into a fixed number of parts.
//
//   - Chunk: fixed number of parts, variable part size.
//   - Split: variable number of parts, fixed part size.
//
// NOTE: If the number of parts exceeds the size of the
// source slice then the resulting slice will contain
// one element slices for each element in the source
// slice and nil slices for the remaining entries.
func Chunk[X any](src []X, numParts int) [][]X {
	srcLen := len(src)
	if srcLen == 0 || numParts < 1 {
		return nil
	}
	if numParts == 1 {
		return [][]X{src}
	}
	chunks := make([][]X, numParts)
	if numParts >= srcLen {
		for i := 0; i < srcLen; i++ {
			chunks[i] = []X{src[i]}
		}
		return chunks
	}
	begin := 0
	for i := 0; i < numParts; i++ {
		end := ((i + 1) * srcLen) / numParts
		chunks[i] = src[begin:end]
		begin = end
	}
	return chunks
}
