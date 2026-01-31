package srv

import (
	"context"
	"net/http"

	"github.com/kw510/zee/pkg/gen/zee/api/v1/apiv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func Init(ctx context.Context) (http.Handler, error) {
	mux := http.NewServeMux()
	path, handler := apiv1connect.NewApiServiceHandler(&ApiServer{})
	mux.Handle(path, handler)

	// Use h2c so we can serve HTTP/2 without TLS.
	srv := h2c.NewHandler(mux, &http2.Server{})

	return srv, nil
}
