package httputil

import (
	"bufio"
	"io"

	"github.com/jamesross3/go-play2/genericsync"
)

type BufioWriterPool interface {
	AcquireWriterFor(io.Writer) *bufio.Writer
	RelinquishWriter(*bufio.Writer)
}

type bufioWriterPool struct {
	pool genericsync.Pool[*bufio.Writer]
}

// AcquireWriterFor wraps the provided writer in a *bufio.Writer.
// The returned writer is ready for use.
func (b *bufioWriterPool) AcquireWriterFor(w io.Writer) *bufio.Writer {
	bufioWriter := b.pool.Get()
	bufioWriter.Reset(w)
	return bufioWriter
}

// RelinquishWriter implements BufioWriterPool.
// The caller should call Flush on the provided writer before calling this method.
func (b *bufioWriterPool) RelinquishWriter(w *bufio.Writer) {
	b.pool.Put(w)
}

func NewBufferedWriterPool(bufSize int) BufioWriterPool {
	return &bufioWriterPool{
		pool: genericsync.NewPool[*bufio.Writer](func() *bufio.Writer {
			return bufio.NewWriterSize(io.Discard, bufSize)
		}),
	}
}
