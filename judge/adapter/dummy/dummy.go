package dummy

import (
	"context"
	"encoding/json"

	"github.com/zhzxdev/azukiiro/client"
	"github.com/zhzxdev/azukiiro/common"
)

type DummyConfig struct {
	Ping string `json:"ping"`
}

type DummyAdapter struct{}

func (d *DummyAdapter) Name() string {
	return "dummy"
}

func (d *DummyAdapter) Judge(ctx context.Context, config common.ProblemConfig, problemData string, solutionData string) error {
	adapterConfig := DummyConfig{}
	json.Unmarshal(config.Judge.Config, &adapterConfig)
	client.PatchSolutionTask(ctx, &client.PatchSolutionTaskRequest{
		Score: 100,
		Metrics: &map[string]float64{
			"cpu": 0,
			"mem": 0,
		},
		Status:  "AC",
		Message: "Well Done! Accepted",
	})
	client.SaveSolutionDetails(ctx, &common.SolutionDetails{
		Jobs: []*common.SolutionDetailsJob{
			{
				Name:       "Group 1",
				Score:      100,
				ScoreScale: 100,
				Status:     "AC",
				Tests: []*common.SolutionDetailsTest{
					{
						Name:    "Test 1",
						Score:   100,
						Status:  "AC",
						Summary: "Accepted",
					},
				},
				Summary: "Accepted",
			},
		},
		Summary: "Accepted\nPing is: `" + adapterConfig.Ping + "`",
	})
	return nil
}
