package httputil

import (
	"bufio"
	"io"
	"net/http"
)

type WrappedResponseWriter interface {
	http.ResponseWriter
	// Done cleans up any resources associated with the WrappedResponseWriter implementation/
	// The caller must call Done() when finished using the WrappedResponseWriter
	Done() error
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	writer io.Writer
}

func (w *wrappedResponseWriter) Write(in []byte) (int, error) {
	return w.writer.Write(in)
}

type ResponseWriterWrapper interface {
	Wrap(rw http.ResponseWriter) WrappedResponseWriter
}

// NewResponseWriterWrapperWithBufio returns a ResponseWriterWrapper which instantiates a new *bufio.Writer
// for each Wrap() call.
func NewResponseWriterWrapperWithBufio(bufSize int) ResponseWriterWrapper {
	wrapper := responseWriterWrapperWithBufio(bufSize)
	return &wrapper
}

type responseWriterWrapperWithBufio int

func (r *responseWriterWrapperWithBufio) Wrap(rw http.ResponseWriter) WrappedResponseWriter {
	bufioWriter := bufio.NewWriterSize(rw, int(*r))
	return &withDoner{
		ResponseWriter: &wrappedResponseWriter{
			writer:         bufioWriter,
			ResponseWriter: rw,
		},
		donerFn: bufioWriter.Flush,
	}
}

// NewResponseWriterWrapperWithBufioWriterPool returns a ResponseWriterWrapper which wraps writer
func NewResponseWriterWrapperWithBufioWriterPool(bufioWriterPool BufioWriterPool) ResponseWriterWrapper {
	return &responseWriterWrapperWithBufioWriterPool{
		bufioWriterPool: bufioWriterPool,
	}
}

type responseWriterWrapperWithBufioWriterPool struct {
	bufioWriterPool BufioWriterPool
}

type withDoner struct {
	http.ResponseWriter
	donerFn
}

type donerFn func() error

func (fn donerFn) Done() error {
	return fn()
}

func (r *responseWriterWrapperWithBufioWriterPool) Wrap(rw http.ResponseWriter) WrappedResponseWriter {
	bufioWriter := r.bufioWriterPool.AcquireWriterFor(rw)
	return &withDoner{
		ResponseWriter: &wrappedResponseWriter{
			writer:         bufioWriter,
			ResponseWriter: rw,
		},
		donerFn: func() error {
			defer r.bufioWriterPool.RelinquishWriter(bufioWriter)
			return bufioWriter.Flush()
		},
	}
}
