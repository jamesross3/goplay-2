package httputil

import (
	"io"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func runBenchmark(b *testing.B, wrapper ResponseWriterWrapper) {
	testFS := &myFS{rootDir: "testdata"}
	fserve := http.FileServer(testFS)

	time.Sleep(100 * time.Millisecond)
	go http.ListenAndServe("localhost:8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrappedResponseWriter := wrapper.Wrap(w)
		defer wrappedResponseWriter.Done()
		fserve.ServeHTTP(wrappedResponseWriter, r)
	}))

	time.Sleep(100 * time.Millisecond)
	b.ResetTimer()

	var totalSize int64
	for i := 0; i < b.N; i++ {
		size, _ := doReq(b)
		totalSize += size
	}
	b.ReportAllocs()
	b.SetBytes(totalSize)
}

func BenchmarkBaseline(b *testing.B) {
	runBenchmark(b, identityResponseWriterWrapper{})
}

func BenchmarkBufio0(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufio(4096))
}

type identityResponseWriterWrapper struct{}

func (i identityResponseWriterWrapper) Wrap(rw http.ResponseWriter) WrappedResponseWriter {
	return withDoner{
		ResponseWriter: rw,
		donerFn:        func() error { return nil },
	}
}

func doReq(t testing.TB) (size int64, duration time.Duration) {
	resp, err := http.Get("http://localhost:8080/30MB.txt")
	require.NoError(t, err)
	defer resp.Body.Close()
	start := time.Now()
	n, err := io.Copy(io.Discard, resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, resp.ContentLength, n)
	return resp.ContentLength, time.Since(start)
}

type myFS struct {
	rootDir string
}

func (m *myFS) Open(name string) (http.File, error) {
	return os.Open(path.Join(m.rootDir, name))
}
