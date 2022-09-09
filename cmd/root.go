package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/toyozaki/connect-go-example/api"
	"github.com/toyozaki/connect-go-example/cmd/utils"
)

var (
	cfgFile string
	appName = "connect-go-example"
	rootCmd = &cobra.Command{
		Use:   appName,
		Short: fmt.Sprintf("%s is a template\n", appName),
		Run: func(cmd *cobra.Command, args []string) {
			addr := utils.GetAddr()
			api.Run(addr)
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is $HOME/.%s.yaml)", appName))

	host := "host"
	rootCmd.Flags().StringP(host, "", "localhost", "host")
	viper.BindPFlag(host, rootCmd.Flags().Lookup(host))

	port := "port"
	rootCmd.Flags().IntP(port, "", 8080, "port")
	viper.BindPFlag(port, rootCmd.Flags().Lookup(port))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigName("." + appName)
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
