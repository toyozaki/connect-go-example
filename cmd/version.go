package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version    = "X.X.X"
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version shows this app's version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("version:", version)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
