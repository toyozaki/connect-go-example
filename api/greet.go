package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bufbuild/connect-go"
	greetv1 "github.com/toyozaki/connect-go-example/gen/greet/v1"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/durationpb"
)

type GreetServer struct{}

func (s *GreetServer) Greet(ctx context.Context, req *connect.Request[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {
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

	res := connect.NewResponse(&greetv1.GreetResponse{
		Greeting: greeting,
	})
	res.Header().Set("Greet-Version", "v1")
	res.Trailer().Set("Greet-Version", "v1")

	// no-ASCII値をheaderで送る場合、base64エンコードが必要になる
	// また、-Binをサフィックスに付ける必要がある
	res.Header().Set(
		"Greet-Emoji-Bin",
		connect.EncodeBinaryHeader([]byte("👋")),
	)

	return res, nil
}

func validateGreetRequest(msg *greetv1.GreetRequest) error {
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

func doGreetWork(ctx context.Context, msg *greetv1.GreetRequest) (string, error) {
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
