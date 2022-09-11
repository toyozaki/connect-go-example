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
			// get headers
			log.Println("Greet-Version in Header", connectErr.Meta().Get("Greet-Version"))
			log.Println("Greet-Version in Trailer", connectErr.Meta().Get("Greet-Version"))

			if retryInfoErr, ok := extractRetryInfo(connectErr); ok {
				log.Fatalln(connectErr.Message(), retryInfoErr.GetRetryDelay())
			}
			log.Fatalln(connectErr.Message(), connectErr.Details())
		}
		log.Fatalln(connect.CodeOf(err))
	}

	log.Println("Greet-Version in Header", res.Header().Get("Greet-Version"))
	log.Println("Greet-Version in Trailer", res.Trailer().Get("Greet-Version"))
	encodedEmoji := res.Header().Get("Greet-Emoji-Bin")
	if decodedEmoji, err := connect.DecodeBinaryHeader(encodedEmoji); err == nil {
		log.Println("Greet-Emoji", string(decodedEmoji))
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
