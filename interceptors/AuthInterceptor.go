package interceptors

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
)

const tokenHeader = "Greet-Token"

var errNoToken = errors.New("no token provided")

type authIntercetor struct{}

func NewAuthInterceptor() *authIntercetor {
	return &authIntercetor{}
}

func (i *authIntercetor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(
		func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				req.Header().Set(tokenHeader, "sample")
			} else if req.Header().Get(tokenHeader) == "" {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					errNoToken,
				)
			}
			return next(ctx, req)
		},
	)
}

func (*authIntercetor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return connect.StreamingClientFunc(
		func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
			conn := next(ctx, spec)
			conn.RequestHeader().Set(tokenHeader, "sample")
			return conn
		})
}

func (i *authIntercetor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return connect.StreamingHandlerFunc(
		func(ctx context.Context, conn connect.StreamingHandlerConn) error {
			if conn.RequestHeader().Get(tokenHeader) == "" {
				return connect.NewError(connect.CodeUnauthenticated, errNoToken)
			}
			return next(ctx, conn)
		})
}
