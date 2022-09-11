package api

import (
	"net/http"

	"github.com/toyozaki/connect-go-example/gen/greet/v1/greetv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func Run(addr string) {
	api := http.NewServeMux()
	api.Handle(greetv1connect.NewGreetServiceHandler(&GreetServer{}))

	mux := http.NewServeMux()
	mux.Handle("/grpc/", http.StripPrefix("/grpc", api))
	http.ListenAndServe(
		addr,
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
