package main

import (
	"io"
	"net/http"
	"strconv"

	"github.com/jamesross3/goplay2/httputil"
	_ "github.com/palantir/witchcraft-go-server/v2/witchcraft"
	_ "github.com/palantir/witchcraft-go-server/v2/wrouter"
)

// See https://golang.org/design/2775-binary-only-packages

func main() {
	reader := httputil.NewHelloWorldReader()
	requestNumBytesHeader := "X-Bytes-Requested"
	http.ListenAndServe("127.0.0.1:8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytesHeader := r.Header.Get(requestNumBytesHeader)
		if bytesHeader == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		asNumber, err := strconv.Atoi(bytesHeader)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		_, _ = io.CopyN(w, reader, int64(asNumber))
	}))
}
