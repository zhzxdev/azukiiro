package judge

import (
	"context"
	"fmt"
	"log"

	"github.com/zhzxdev/azukiiro/client"
	"github.com/zhzxdev/azukiiro/common"
	"github.com/zhzxdev/azukiiro/judge/adapter"
	"github.com/zhzxdev/azukiiro/storage"
)

func judge(ctx context.Context, res *client.PollSolutionResponse) error {
	problemData, err := storage.PrepareFile(ctx, res.ProblemDataUrl, res.ProblemDataHash)
	if err != nil {
		return err
	}
	solutionData, err := storage.PrepareFile(ctx, res.SolutionDataUrl, res.SolutionDataHash)
	if err != nil {
		return err
	}
	adapter, ok := adapter.Get(res.ProblemConfig.Judge.Adapter)
	if !ok {
		return client.PatchSolutionTask(ctx, &client.PatchSolutionTaskRequest{
			Score:   0,
			Status:  "Error",
			Message: "Judge adapter not found",
		})
	}
	return adapter.Judge(ctx, res.ProblemConfig, problemData, solutionData)
}

func Poll(ctx context.Context) (bool, error) {
	res, err := client.PollSolution(ctx, &client.PollSolutionRequest{})
	if err != nil {
		return false, err
	}

	if res.TaskId == "" {
		// No pending tasks
		return false, nil
	}

	ctx = client.WithSolutionTask(ctx, res.SolutionId, res.TaskId)

	if res.ErrMsg != "" {
		// Server side error occurred
		client.PatchSolutionTask(ctx, &client.PatchSolutionTaskRequest{
			Score:   0,
			Status:  "Error",
			Message: "Server side error occurred",
		})
		client.CompleteSolutionTask(ctx)
		return true, nil
	}

	log.Println("Got task:", res.TaskId)
	log.Println("SolutionId:", res.SolutionId)

	err = judge(ctx, res)
	if err != nil {
		log.Println("Judge finished with error:", err)
		err = client.SaveSolutionDetails(ctx, &common.SolutionDetails{
			Jobs:    []*common.SolutionDetailsJob{},
			Summary: fmt.Sprintf("An Error has occurred:\n\n```\n%s\n```", err),
		})
		if err != nil {
			log.Println("Save details failed:", err)
		}
		err = client.PatchSolutionTask(ctx, &client.PatchSolutionTaskRequest{
			Score:   0,
			Status:  "Error",
			Message: "Judge error",
		})
		if err != nil {
			log.Println("Patch task failed:", err)
		}
	} else {
		log.Println("Judge finished")
	}
	err = client.CompleteSolutionTask(ctx)
	if err != nil {
		log.Println("Complete task failed:", err)
	}

	return true, nil
}
