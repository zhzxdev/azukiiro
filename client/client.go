package client

import (
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

var http = resty.New()

func GetDefaultHTTPClient() *resty.Client {
	return http
}

func InitFromConfig() {
	serverAddr := viper.GetString("serverAddr")
	if serverAddr == "" {
		log.Fatalln("Server address not set")
	}
	http.SetBaseURL(serverAddr)
	runnerId := viper.GetString("runnerId")
	if runnerId == "" {
		log.Fatalln("Runner ID not set")
	}
	runnerKey := viper.GetString("runnerKey")
	if runnerKey == "" {
		log.Fatalln("Runner key not set")
	}
	http.SetHeaders(map[string]string{
		"X-AOI-Runner-Id":  runnerId,
		"X-AOI-Runner-Key": runnerKey,
	})
}
