package httputil

import (
	"context"
	"io"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"github.com/jamesross3/goplay2/resources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	respWriterWrapperPointer *ResponseWriterWrapper
)

// TestMain sets up a file server with a respWriterWrapperPointer for individual benchmark tests to modify when they start running.
// for now, it must be run without any parallelism (`go test -p=1 -bench=.`)
func TestMain(m *testing.M) {
	testFS := &myFS{rootDir: "testdata"}
	fserve := http.FileServer(testFS)

	time.Sleep(100 * time.Millisecond)
	go http.ListenAndServe("localhost:8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrappedResponseWriter := (*respWriterWrapperPointer).Wrap(w)
		defer wrappedResponseWriter.Done()
		fserve.ServeHTTP(wrappedResponseWriter, r)
	}))
	time.Sleep(time.Second)
	os.Exit(m.Run())
}

func runBenchmark(b *testing.B, wrapper ResponseWriterWrapper) {
	respWriterWrapperPointer = &wrapper
	for i := 0; i < b.N; i++ {
		size, _ := doReq(b)
		b.SetBytes(size)
	}
	b.ReportAllocs()
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

type identityResponseWriterWrapper struct{}

func (i identityResponseWriterWrapper) Wrap(rw http.ResponseWriter) WrappedResponseWriter {
	return withDoner{
		ResponseWriter: rw,
		donerFn:        func() error { return nil },
	}
}

func BenchmarkBaseline(b *testing.B) { // 178           6177217 ns/op        5078.46 MB/s       39182 B/op         81 allocs/op
	runBenchmark(b, identityResponseWriterWrapper{})
}

func BenchmarkBufio0(b *testing.B) { // 174           6391928 ns/op           43499 B/op         86 allocs/op
	runBenchmark(b, NewResponseWriterWrapperWithBufio(4096))
}

func BenchmarkBufio1(b *testing.B) { // 178           6307001 ns/op           47460 B/op         86 allocs/op
	runBenchmark(b, NewResponseWriterWrapperWithBufio(4096<<1))
}

func BenchmarkBufio2(b *testing.B) { // 174           6288829 ns/op           55968 B/op         86 allocs/op
	runBenchmark(b, NewResponseWriterWrapperWithBufio(4096<<2))
}

func BenchmarkBufio3(b *testing.B) { // 153           7014077 ns/op           72604 B/op         86 allocs/op
	runBenchmark(b, NewResponseWriterWrapperWithBufio(4096<<3))
}

func BenchmarkBufio4(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufio(4096<<4))
}

func BenchmarkBufio5(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufio(4096<<5))
}

func BenchmarkBufio6(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufio(4096<<6))
}

func BenchmarkBufio7(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufio(4096<<7))
}

func BenchmarkBufio8(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufio(4096<<8))
}

func BenchmarkBufio9(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufio(4096<<9))
}

func BenchmarkBufio10(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufio(4096<<10))
}

// pool

func BenchmarkBufioPool0(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufioWriterPool(NewBufferedWriterPool(4096)))
}

func BenchmarkBufioPool1(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufioWriterPool(NewBufferedWriterPool(4096<<1)))
}

func BenchmarkBufioPool2(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufioWriterPool(NewBufferedWriterPool(4096<<2)))
}

func BenchmarkBufioPool3(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufioWriterPool(NewBufferedWriterPool(4096<<3)))
}

func BenchmarkBufioPool4(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufioWriterPool(NewBufferedWriterPool(4096<<4)))
}

func BenchmarkBufioPool5(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufioWriterPool(NewBufferedWriterPool(4096<<5)))
}

func BenchmarkBufioPool6(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufioWriterPool(NewBufferedWriterPool(4096<<6)))
}

func BenchmarkBufioPool7(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufioWriterPool(NewBufferedWriterPool(4096<<7)))
}

func BenchmarkBufioPool8(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufioWriterPool(NewBufferedWriterPool(4096<<8)))
}

func BenchmarkBufioPool9(b *testing.B) { // 199           5252033 ns/op        5973.06 MB/s       76431 B/op         83 alloc
	runBenchmark(b, NewResponseWriterWrapperWithBufioWriterPool(NewBufferedWriterPool(4096<<9)))
	_ = resources.Blocks(context.TODO())
}

func BenchmarkBufioPool10(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufioWriterPool(NewBufferedWriterPool(4096<<10)))
}

func BenchmarkBufioPool11(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufioWriterPool(NewBufferedWriterPool(4096<<11)))
}

func BenchmarkBufioPool12(b *testing.B) {
	runBenchmark(b, NewResponseWriterWrapperWithBufioWriterPool(NewBufferedWriterPool(4096<<12)))
}

func TestLen(t *testing.T) {
	str := "hello world!éº0˚∂∂∂∂∂∂∂"
	b := []byte(str)
	assert.Equal(t, len(str), len(b))
}
