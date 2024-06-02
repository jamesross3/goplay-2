package genericsync

import "sync"

type Pool[T any] interface {
	Get() T
	Put(T)
}

type pool[T any] struct {
	pool *sync.Pool
}

// AcquireFor implements BufferPool.
func (p *pool[T]) Get() T {
	return p.pool.Get().(T)
}

func (p *pool[T]) Put(t T) {
	p.pool.Put(t)
}

func NewPool[T any](create func() T) Pool[T] {
	return &pool[T]{
		pool: &sync.Pool{
			New: func() any {
				return create()
			},
		},
	}
}
