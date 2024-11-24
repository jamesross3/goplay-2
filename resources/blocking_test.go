package resources

import (
	"context"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlocks(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go someWork(wg)
	assert.NoError(t, Blocks(context.TODO()))
	wg.Wait()
}

func someWork(wg *sync.WaitGroup) error {
	defer wg.Done()
	resp, err := http.Get("https://example.com")
	if err != nil {
		return err
	}
	if resp.Body == nil {
		return nil
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	return err
}
