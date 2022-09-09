package client

import (
	"context"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	greetv1 "github.com/toyozaki/connect-go-example/gen/greet/v1"
	"github.com/toyozaki/connect-go-example/gen/greet/v1/greetv1connect"
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
		log.Fatalln(err)
	}
	log.Println(res.Msg.Greeting)
}
