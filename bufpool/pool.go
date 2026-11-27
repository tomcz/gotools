package bufpool

import (
	"bytes"
	"sync"
)

// Pool implements a garbage-collection-friendly pool of output buffers backed by [sync.Pool].
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

// Borrow a buffer from the pool.
//
// Each borrowed buffer must be returned back to the pool to allow reuse and prevent memory leaks.
func (p *Pool) Borrow() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

// Return borrowed buffer back into the pool.
func (p *Pool) Return(b *bytes.Buffer) {
	b.Reset()
	p.pool.Put(b)
}
