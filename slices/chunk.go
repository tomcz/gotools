package slices

// Chunk given slice into a fixed number of parts.
func Chunk[X any](src []X, numParts int) [][]X {
	srcLen := len(src)
	if srcLen == 0 || numParts < 1 {
		return nil
	}
	if numParts == 1 {
		return [][]X{src}
	}
	chunks := make([][]X, numParts)
	if numParts > srcLen {
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
