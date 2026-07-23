package html

import (
	"bytes"
	"sync"
)

var bufPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func bufBorrow() *bytes.Buffer {
	return bufPool.Get().(*bytes.Buffer)
}

func bufReturn(b *bytes.Buffer) {
	b.Reset()
	bufPool.Put(b)
}
