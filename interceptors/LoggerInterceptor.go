package interceptors

import (
	"context"
	"log"

	"github.com/bufbuild/connect-go"
)

func NewLoggerInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(
			func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
				log.Println(req.Spec().Procedure)
				return next(ctx, req)
			},
		)
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
