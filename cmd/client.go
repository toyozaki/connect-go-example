package cmd

import (
	"github.com/spf13/cobra"
	"github.com/toyozaki/connect-go-example/client"
)

type GreetSubCommand string

var (
	UNARY_GREET         GreetSubCommand = "unary"
	CLIENT_STREAM_GREET GreetSubCommand = "c-stream"
)

var (
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "greet client",
	}
	clientUnaryGreetCmd = &cobra.Command{
		Use:   string(UNARY_GREET),
		Short: "unary greet",
		Run: func(cmd *cobra.Command, args []string) {
			greet(cmd, UNARY_GREET)
		},
	}
	clientClientStreamGreetCmd = &cobra.Command{
		Use:   string(CLIENT_STREAM_GREET),
		Short: "client streaming greet",
		Run: func(cmd *cobra.Command, args []string) {
			greet(cmd, CLIENT_STREAM_GREET)
		},
	}
)

func greet(cmd *cobra.Command, subcommand GreetSubCommand) {
	name, _ := cmd.Flags().GetString("name")
	addr, _ := cmd.Flags().GetString("addr")
	client := client.NewGreetClient(addr)
	switch subcommand {
	case UNARY_GREET:
		client.UnaryGreet(name)
	case CLIENT_STREAM_GREET:
		client.ClientStreamGreet(name)
	}
}

func init() {
	clientCmd.PersistentFlags().String("addr", "http://localhost:8080/grpc/", "addr")
	clientCmd.PersistentFlags().String("name", "toyozaki", "your name")
	clientCmd.AddCommand(clientUnaryGreetCmd)
	clientCmd.AddCommand(clientClientStreamGreetCmd)

	rootCmd.AddCommand(clientCmd)
}
