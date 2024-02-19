package main

import (
	"context"
)

type debugstr string

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ctx1 := context.WithValue(ctx, debugstr("debug"), "please")
	ctx2, _ := context.WithCancel(ctx1)
	cancel()
	<-ctx2.Done()
}
