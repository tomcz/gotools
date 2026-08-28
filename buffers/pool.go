package buffers

import (
	"bytes"
	"sync"
)

// Pool implements a garbage-collection-friendly pool of output buffers backed by `sync.Pool`.
//
// Each borrowed buffer must be returned back to the pool to allow reuse and prevent memory leaks.
type Pool struct {
	pool sync.Pool
}

func New() *Pool {
	return &Pool{
		pool: sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
}

// Borrow buffer from the pool.
func (p *Pool) Borrow() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

// Return borrowed buffer back into the pool.
func (p *Pool) Return(b *bytes.Buffer) {
	b.Reset()
	p.pool.Put(b)
}
