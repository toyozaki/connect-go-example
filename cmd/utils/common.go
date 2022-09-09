package utils

import (
	"strconv"

	"github.com/spf13/viper"
)

func GetAddr() string {
	host := viper.GetString("host")
	port := viper.GetInt("port")
	return host + ":" + strconv.Itoa(port)
}
