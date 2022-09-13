package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bufbuild/connect-go"
	greetv1 "github.com/toyozaki/connect-go-example/gen/greet/v1"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/durationpb"
)

type GreetServer struct{}

func (s *GreetServer) UnaryGreet(ctx context.Context, req *connect.Request[greetv1.UnaryGreetRequest]) (*connect.Response[greetv1.UnaryGreetResponse], error) {
	log.Println("Request headers: ", req.Header())

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := validateGreetRequest(req.Msg); err != nil {
		return nil, err
	}

	greeting, err := doGreetWork(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	res := connect.NewResponse(&greetv1.UnaryGreetResponse{
		Greeting: greeting,
	})

	setCommonHeader(res)
	setCommonTrailer(res)

	return res, nil
}

func (*GreetServer) ClientStreamGreet(ctx context.Context, stream *connect.ClientStream[greetv1.ClientStreamGreetRequest]) (*connect.Response[greetv1.ClientStreamGreetResponse], error) {
	log.Println("Request headers: ", stream.RequestHeader())
	var greeting strings.Builder

	for stream.Receive() {
		g := fmt.Sprintf("Hello, %s!\n", stream.Msg().Name)
		if _, err := greeting.WriteString(g); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}
	if err := stream.Err(); err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}
	res := connect.NewResponse(&greetv1.ClientStreamGreetResponse{
		Greeting: greeting.String(),
	})

	setCommonHeader(res)
	setCommonTrailer(res)

	return res, nil

}

func validateGreetRequest(msg *greetv1.UnaryGreetRequest) error {
	if msg.Name == "invalid" {
		connectErr := connect.NewError(connect.CodeInvalidArgument, errors.New("invalid name"))
		connectErr = setGreetVersionToErr(connectErr)
		return connectErr
	}
	if msg.Name == "trasient" {
		return newTransientError()
	}
	return nil
}

func doGreetWork(ctx context.Context, msg *greetv1.UnaryGreetRequest) (string, error) {
	return fmt.Sprintf("Hello, %s!", msg.Name), nil
}

func newTransientError() error {
	err := connect.NewError(
		connect.CodeUnavailable,
		errors.New("overloaded: back off and retry"),
	)
	retryInfo := &errdetails.RetryInfo{
		RetryDelay: durationpb.New(10 * time.Second),
	}
	if detail, detailErr := connect.NewErrorDetail(retryInfo); detailErr == nil {
		err.AddDetail(detail)
	}

	err = setGreetVersionToErr(err)

	return err
}

func setGreetVersionToErr(connectErr *connect.Error) *connect.Error {
	connectErr.Meta().Set("Greet-Version", "v1")
	return connectErr
}

func setCommonHeader[T any](res *connect.Response[T]) {
	res.Header().Set("Greet-Version", "v1")

	// no-ASCII値をheaderで送る場合、base64エンコードが必要になる
	// また、-Binをサフィックスに付ける必要がある
	res.Header().Set(
		"Greet-Emoji-Bin",
		connect.EncodeBinaryHeader([]byte("👋")),
	)
}

func setCommonTrailer[T any](res *connect.Response[T]) {
	res.Trailer().Set("Greet-Version", "v1")
}
