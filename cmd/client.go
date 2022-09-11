package cmd

import (
	"github.com/spf13/cobra"
	"github.com/toyozaki/connect-go-example/client"
)

var (
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "greet client",
	}
	clientGreetCmd = &cobra.Command{
		Use:   "greet",
		Short: "greet!",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			addr, _ := cmd.Flags().GetString("addr")
			client := client.NewGreetClient(addr)
			client.Greet(name)
		},
	}
)

func init() {
	clientCmd.PersistentFlags().String("addr", "http://localhost:8080/grpc/", "addr")
	clientCmd.PersistentFlags().String("name", "toyozaki", "your name")
	clientCmd.AddCommand(clientGreetCmd)
	rootCmd.AddCommand(clientCmd)
}
