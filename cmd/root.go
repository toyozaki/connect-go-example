package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	appName = "connect-go-example"
	rootCmd = &cobra.Command{
		Use:   appName,
		Short: fmt.Sprintf("%s is a template\n", appName),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s is a template\n", appName)
			fmt.Println("host:", viper.GetString("host"))
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
	rootCmd.Flags().StringP(host, "", "", "host")
	viper.BindPFlag(host, rootCmd.Flags().Lookup(host))
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
