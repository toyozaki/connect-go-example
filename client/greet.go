package client

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	greetv1 "github.com/toyozaki/connect-go-example/gen/greet/v1"
	"github.com/toyozaki/connect-go-example/gen/greet/v1/greetv1connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type Client struct {
	greetv1connect.GreetServiceClient
}

func NewGreetClient(addr string) *Client {
	client := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		addr,
		// default: using connect protocol
		// check server log's content-type
		// connect.WithGRPC(),
		connect.WithGRPCWeb(),
	)
	return &Client{
		client,
	}
}

func (c *Client) Greet(name string) {
	res, err := c.GreetServiceClient.Greet(
		context.Background(),
		connect.NewRequest(&greetv1.GreetRequest{Name: name}),
	)
	if err != nil {
		if connectErr := new(connect.Error); errors.As(err, &connectErr) {
			if retryInfoErr, ok := extractRetryInfo(connectErr); ok {
				log.Fatalln(connectErr.Message(), retryInfoErr.GetRetryDelay())
			}
			log.Fatalln(connectErr.Message(), connectErr.Details())
		}
		log.Fatalln(connect.CodeOf(err))
	}
	log.Println(res.Msg.Greeting)
}

func extractRetryInfo(connectErr *connect.Error) (*errdetails.RetryInfo, bool) {
	for _, detail := range connectErr.Details() {
		msg, valueErr := detail.Value()
		if valueErr != nil {
			continue
		}
		if retryInfo, ok := msg.(*errdetails.RetryInfo); ok {
			return retryInfo, true
		}
	}
	return nil, false
}
