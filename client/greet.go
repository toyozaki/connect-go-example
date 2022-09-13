package client

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/bufbuild/connect-go"
	greetv1 "github.com/toyozaki/connect-go-example/gen/greet/v1"
	"github.com/toyozaki/connect-go-example/gen/greet/v1/greetv1connect"
	"github.com/toyozaki/connect-go-example/interceptors"
	"golang.org/x/net/http2"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type Client struct {
	greetv1connect.GreetServiceClient
}

func NewGreetClient(addr string) *Client {
	interceptors := connect.WithInterceptors(interceptors.NewAuthInterceptor())

	httpClient := &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
				// If you're also using this client for non-h2c traffic, you may want
				// to delegate to tls.Dial if the network isn't TCP or the addr isn't
				// in an allowlist.
				return net.Dial(network, addr)
			},
			// Don't forget timeouts!
		},
	}

	client := greetv1connect.NewGreetServiceClient(
		httpClient,
		addr,
		// default: using connect protocol
		// check server log's content-type
		// connect.WithGRPC(),
		connect.WithGRPCWeb(),
		interceptors,
	)
	return &Client{
		client,
	}
}

func (c *Client) UnaryGreet(name string) {
	res, err := c.GreetServiceClient.UnaryGreet(
		context.Background(),
		connect.NewRequest(&greetv1.UnaryGreetRequest{Name: name}),
	)

	failOnError(err)

	printGreetResponse(res)
}

func (c *Client) ClientStreamGreet(name string) {
	greet := c.GreetServiceClient.ClientStreamGreet(context.Background())
	var res *connect.Response[greetv1.ClientStreamGreetResponse]
	for i := 0; i < 3; i++ {
		err := greet.Send(&greetv1.ClientStreamGreetRequest{
			Name: strconv.Itoa(i) + " : " + name,
		})

		if err != nil {
			break
		}
	}

	res, err := greet.CloseAndReceive()
	failOnError(err)

	printGreetResponse(res)
}

func failOnError(err error) {
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

func printGreetResponse[T greetv1.UnaryGreetResponse | greetv1.ClientStreamGreetResponse](res *connect.Response[T]) {
	log.Println("Greet-Version in Header", res.Header().Get("Greet-Version"))
	log.Println("Greet-Version in Trailer", res.Trailer().Get("Greet-Version"))
	encodedEmoji := res.Header().Get("Greet-Emoji-Bin")
	if decodedEmoji, err := connect.DecodeBinaryHeader(encodedEmoji); err == nil {
		log.Println("Greet-Emoji", string(decodedEmoji))
	}

	log.Println(res.Msg)
}
