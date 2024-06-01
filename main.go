package main

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-server/v2/config"
	"github.com/palantir/witchcraft-go-server/v2/witchcraft"
	"github.com/palantir/witchcraft-go-server/v2/wrouter"
)

// See https://golang.org/design/2775-binary-only-packages

func main() {
	witchcraft.
		NewServer().
		WithInstallConfig(config.Install{
			Server: config.Server{
				Address:        "0.0.0.0",
				Port:           8443,
				ManagementPort: 8444,
				ContextPath:    "/jross",
			},
			UseConsoleLog: true,
		}).
		WithRuntimeConfig(config.Runtime{
			LoggerConfig: &config.LoggerConfig{
				Level: wlog.DebugLevel,
			},
		}).
		WithClientAuth(tls.NoClientCert).
		WithSelfSignedCertificate().
		WithInitFunc(initFunc).
		Start()

}

func initFunc(ctx context.Context, initInfo witchcraft.InitInfo) (func(), error) {
	initInfo.Router.Get("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}), wrouter.RouteMiddleware(func(rw http.ResponseWriter, r *http.Request, reqVals wrouter.RequestVals, next wrouter.RouteRequestHandler) {
		ctx, _ := context.WithCancelCause(r.Context())
		req := r.WithContext(ctx)
		next(rw, req, reqVals)
	}))
	return nil, nil
}
