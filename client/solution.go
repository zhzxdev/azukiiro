package client

import (
	"context"
	"encoding/json"

	"github.com/zhzxdev/azukiiro/azukiiro/common"
	"github.com/zhzxdev/azukiiro/azukiiro/storage"
)

type PatchSolutionTaskRequest struct {
	Score   int                 `json:"score"`
	Metrics *map[string]float64 `json:"metrics,omitempty"`
	Status  string              `json:"status"`
	Message string              `json:"message"`
}

// Context helpers
type taskContextKey int

const (
	solutionIdKey taskContextKey = iota
	taskIdKey     taskContextKey = iota
)

func WithSolutionTask(ctx context.Context, solutionId string, taskId string) context.Context {
	ctx = context.WithValue(ctx, solutionIdKey, solutionId)
	ctx = context.WithValue(ctx, taskIdKey, taskId)
	return ctx
}

func LoadSolutionTask(ctx context.Context) (solutionId string, taskId string) {
	solutionId = ctx.Value(solutionIdKey).(string)
	taskId = ctx.Value(taskIdKey).(string)
	return
}

type PollSolutionRequest struct {
}

type PollSolutionResponse struct {
	TaskId           string               `json:"taskId"`
	SolutionId       string               `json:"solutionId"`
	ProblemConfig    common.ProblemConfig `json:"problemConfig"`
	ProblemDataUrl   string               `json:"problemDataUrl"`
	ProblemDataHash  string               `json:"problemDataHash"`
	SolutionDataUrl  string               `json:"solutionDataUrl"`
	SolutionDataHash string               `json:"solutionDataHash"`
	ErrMsg           string               `json:"errMsg"`
}

func PollSolution(ctx context.Context, req *PollSolutionRequest) (*PollSolutionResponse, error) {
	res := &PollSolutionResponse{}
	raw, err := http.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(res).
		Post("/api/runner/solution/poll")
	err = loadError(raw, err)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SaveSolutionDetails(ctx context.Context, details *common.SolutionDetails) error {
	url, err := GetSolutionTaskDetailsUrl(ctx, "upload")
	if err != nil {
		return err
	}
	str, err := json.Marshal(details)
	if err != nil {
		return err
	}
	return storage.Upload(ctx, url, str)
}

func PatchSolutionTask(ctx context.Context, req *PatchSolutionTaskRequest) error {
	solutionId, taskId := LoadSolutionTask(ctx)
	raw, err := http.R().
		SetContext(ctx).
		SetBody(req).
		Patch("/api/runner/solution/task/" + solutionId + "/" + taskId)
	return loadError(raw, err)
}

func CompleteSolutionTask(ctx context.Context) error {
	solutionId, taskId := LoadSolutionTask(ctx)
	raw, err := http.R().
		SetContext(ctx).
		Post("/api/runner/solution/task/" + solutionId + "/" + taskId + "/complete")
	return loadError(raw, err)
}

type UrlResponse struct {
	Url string `json:"url"`
}

func GetSolutionTaskDetailsUrl(ctx context.Context, urlType string) (string, error) {
	solutionId, taskId := LoadSolutionTask(ctx)
	res := &UrlResponse{}
	raw, err := http.R().
		SetContext(ctx).
		SetResult(res).
		Get("/api/runner/solution/task/" + solutionId + "/" + taskId + "/details/" + urlType)
	err = loadError(raw, err)
	if err != nil {
		return "", err
	}
	return res.Url, nil
}
