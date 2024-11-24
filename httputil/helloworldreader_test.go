package httputil_test

import (
	"testing"

	"github.com/jamesross3/go-play2/httputil"
	"github.com/stretchr/testify/assert"
)

func TestHelloWorldReader(t *testing.T) {
	b := make([]byte, 8)
	h := httputil.NewHelloWorldReader()
	n, err := h.Read(b)
	assert.NoError(t, err)
	assert.Len(t, b, n)
	assert.Equal(t, " Hello w", string(b))

	b = make([]byte, 30)
	n, err = h.Read(b)
	assert.NoError(t, err)
	assert.Len(t, b, n)
}

// var toWrite = []byte(` Hello world!`)

// type helloWorldReader struct{}

// func (h helloWorldReader) Read(b []byte) (int, error) {
// 	var written int
// 	var oneWriteLength int = len(toWrite)
// 	var lenB = len(b)
// 	for written < lenB {
// 		numToWrite := oneWriteLength
// 		leftToWrite := lenB - written
// 		if leftToWrite < oneWriteLength {
// 			numToWrite = leftToWrite
// 		}
// 		copy(b[written:], toWrite[:numToWrite])
// 		written += numToWrite
// 		if leftToWrite < oneWriteLength {
// 			break
// 		}
// 	}
// 	return written, nil
// }

// func NewHelloWorldReader() io.Reader {
// 	return helloWorldReader{}
// }
