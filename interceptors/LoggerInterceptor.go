package interceptors

import (
	"context"
	"errors"
	"log"

	"github.com/bufbuild/connect-go"
)

const tokenHeader = "Greet-Token"

func NewAuthInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(
			func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
				if req.Spec().IsClient {
					req.Header().Set(tokenHeader, "sample")
				} else if req.Header().Get(tokenHeader) == "" {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						errors.New("no token provided"),
					)
				}

				return next(ctx, req)
			},
		)
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

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
