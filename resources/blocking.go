package resources

import (
	"context"
	"os"
	"runtime/pprof"
)

func Blocks(ctx context.Context) error {
	blockProfile := pprof.Lookup("block")
	if blockProfile == nil {
		return nil
	}
	blockProfile.WriteTo(os.Stdout, 1)
	return nil
}
