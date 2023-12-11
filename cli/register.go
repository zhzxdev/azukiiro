package cli

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zhzxdev/azukiiro/client"
	"github.com/zhzxdev/azukiiro/common"
)

type registerArgs struct {
	ServerAddr string
	Force      bool
	Name       string
	Labels     string
	Token      string
}

func splitLabels(input string) []string {
	return strings.Split(input, ",")
}

func runRegister(ctx context.Context, regArgs *registerArgs) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		log.Println("Registering runner")
		log.Println("ServerAddr:", regArgs.ServerAddr)

		runnerKey := viper.GetString("runnerKey")
		if runnerKey != "" && !regArgs.Force {
			log.Println("Runner already registered, exiting...")
			return nil
		}

		viper.Set("serverAddr", regArgs.ServerAddr)

		http := client.GetDefaultHTTPClient()
		http.SetBaseURL(regArgs.ServerAddr)

		if regArgs.Name == "" {
			name, err := os.Hostname()
			if err != nil {
				log.Fatalln(err)
			}
			regArgs.Name = name
		}

		if regArgs.Labels == "" {
			regArgs.Labels = "default"
		}

		req := &client.RegisterRequest{
			Name:              regArgs.Name,
			Version:           common.GetVersion(),
			Labels:            splitLabels(regArgs.Labels),
			RegistrationToken: regArgs.Token,
		}

		res, err := client.Register(ctx, req)

		if err != nil {
			log.Fatalln(err)
		}

		log.Println("RunnerId:", res.RunnerId)
		viper.Set("runnerId", res.RunnerId)
		viper.Set("runnerKey", res.RunnerKey)
		viper.WriteConfig()

		log.Println("Runner registered successfully")

		return nil
	}
}
