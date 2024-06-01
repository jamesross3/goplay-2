package httputil

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBuf(t *testing.T) {
	fmt.Printf("Hello\n")
	testFS := &myFS{rootDir: "testdata"}
	fserve := http.FileServer(testFS)

	bufSizer := new(mysizer)
	go http.ListenAndServe("localhost:8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bufioWriter := bufio.NewWriterSize(w, bufSizer.size())
		defer bufioWriter.Flush()
		wrapped := &wrappedResponseWriter{
			ResponseWriter: w,
			writer:         bufioWriter,
		}
		fserve.ServeHTTP(wrapped, r)
	}))
	time.Sleep(time.Second)
	numReqs := 100
	type Result struct {
		BufSize             int
		BytesPerMicrosecond int64
	}
	var results []Result
	for i := 0; i < 11; i++ {
		bufSize := 4096 << i
		*bufSizer = mysizer(bufSize)
		var totalSize int64
		var totalTime time.Duration
		for i := 0; i < numReqs; i++ {
			size, duration := doReq(t)
			totalSize += size
			totalTime += duration
		}
		results = append(results, Result{
			BufSize:             bufSize,
			BytesPerMicrosecond: totalSize / totalTime.Microseconds(),
		})
		time.Sleep(time.Second)
	}
	asJSON, err := json.MarshalIndent(results, "", "  ")
	require.NoError(t, err)
	t.Logf("\n\n%s\n\n", string(asJSON))
}

type mysizer int

func (m *mysizer) size() int {
	return int(*m)
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	writer io.Writer
}

func (w *wrappedResponseWriter) Write(in []byte) (int, error) {
	return w.writer.Write(in)
}

// no bufio: 3404 bytes/microsecond
// bufio 4096 (golang http default): Avg throughput: 3408 bytes/microsecond (https://stackoverflow.com/questions/26033853/will-http-responsewriters-write-function-buffer-in-go)
// bufio 4096*2: 3288 bytes/microsecond
// bufio 4096*4: 3403 bytes/microsecond

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
