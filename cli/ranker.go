package cli

import (
	"context"
	"log"
	"time"

	"github.com/zhzxdev/azukiiro/client"
	"github.com/zhzxdev/azukiiro/db"
	"github.com/zhzxdev/azukiiro/ranker"

	"github.com/spf13/cobra"
)

type rankerArgs struct {
	pollInterval float32
}

func runRanker(ctx context.Context, rankerArgs *rankerArgs) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		log.Println("Starting ranker")
		ctx, cleanup := db.WithMongo(ctx)
		defer cleanup()

		client.InitFromConfig()
		for {
			cont, err := ranker.Poll(ctx)
			if err != nil {
				log.Println("Error:", err)
			}
			waitDur := time.Duration(0)
			if !cont {
				waitDur = time.Duration(rankerArgs.pollInterval) * time.Second
			}
			timer := time.NewTimer(waitDur)
			select {
			case <-ctx.Done():
				if !timer.Stop() {
					<-timer.C
				}
				return nil
			case <-timer.C:
			}
		}
	}
}
