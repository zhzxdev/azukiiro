package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zhzxdev/azukiiro/storage"
)

type ranklistContextKey int

const injectionKey ranklistContextKey = iota

type ranklistContext struct {
	TaskId    string
	ContestId string
}

func WithRanklistTask(ctx context.Context, taskId, contestId string) context.Context {
	ctx = context.WithValue(ctx, injectionKey, ranklistContext{
		TaskId:    taskId,
		ContestId: contestId,
	})
	return ctx
}

func LoadRanklistCtx(ctx context.Context) (string, string) {
	inj := ctx.Value(injectionKey).(ranklistContext)
	return inj.TaskId, inj.ContestId
}

type RanklistSettings struct {
	ShowAfter  int `json:"showAfter"`
	ShowBefore int `json:"showBefore"`
}

type RanklistDTO struct {
	Key      string           `json:"key"`
	Name     string           `json:"name"`
	Settings RanklistSettings `json:"settings"`
}

type PollRanklistRequest struct {
}

type PollRanklistResponse struct {
	TaskId    string        `json:"taskId"`
	ContestId string        `json:"contestId"`
	Ranklists []RanklistDTO `json:"ranklists"`
}

func PollRanklist(ctx context.Context, req *PollRanklistRequest) (*PollRanklistResponse, error) {
	res := &PollRanklistResponse{}
	raw, err := http.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(res).
		Post("/api/runner/ranklist/poll")
	err = loadError(raw, err)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type RanklistTopstarItemMutation struct {
	Score float64 `json:"score"`
	Ts    int     `json:"ts"`
}

type RanklistTopstarItem struct {
	UserId    string                         `json:"userId"`
	Mutations []*RanklistTopstarItemMutation `json:"mutations"`
}

type RanklistTopstar struct {
	List []*RanklistTopstarItem `json:"list"`
}

type RanklistParticipantColumn struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RanklistParticipantItemColumn struct {
	Content string `json:"content"`
}

type RanklistParticipantItem struct {
	Rank    int                              `json:"rank"`
	UserId  string                           `json:"userId"`
	Columns []*RanklistParticipantItemColumn `json:"columns"`
}

type RanklistParticipant struct {
	Columns []*RanklistParticipantColumn `json:"columns"`
	List    []*RanklistParticipantItem   `json:"list"`
}

type RanklistMetadata struct {
	GeneratedAt int    `json:"generatedAt"`
	Description string `json:"description"`
}

type Ranklist struct {
	Topstar     *RanklistTopstar     `json:"topstar,omitempty"`
	Participant *RanklistParticipant `json:"participant"`
	Metadata    *RanklistMetadata    `json:"metadata"`
}

type GetRanklistUploadUrlsResponse = []struct {
	Key string `json:"key"`
	Url string `json:"url"`
}

func GetRanklistUploadUrls(ctx context.Context) (*GetRanklistUploadUrlsResponse, error) {
	taskId, contestId := LoadRanklistCtx(ctx)
	res := &GetRanklistUploadUrlsResponse{}
	raw, err := http.R().
		SetContext(ctx).
		SetResult(res).
		Get("/api/runner/ranklist/task/" + contestId + "/" + taskId + "/uploadUrls")
	err = loadError(raw, err)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SaveRanklist(ctx context.Context, ranklist map[string]*Ranklist) error {
	res, err := GetRanklistUploadUrls(ctx)
	if err != nil {
		return err
	}
	urls := make(map[string]string)
	for _, value := range *res {
		urls[value.Key] = value.Url
	}
	for key, value := range ranklist {
		str, err := json.Marshal(value)
		if err != nil {
			return err
		}
		err = storage.Upload(ctx, urls[key], str)
		if err != nil {
			return err
		}
	}
	return nil
}

type CompleteRanklistTaskRequest struct {
	RanklistLastSolutionId string `json:"ranklistLastSolutionId"`
}

func CompleteRanklistTask(ctx context.Context, req *CompleteRanklistTaskRequest) error {
	taskId, contestId := LoadRanklistCtx(ctx)
	raw, err := http.R().
		SetContext(ctx).
		SetBody(req).
		Post("/api/runner/ranklist/task/" + contestId + "/" + taskId + "/complete")
	return loadError(raw, err)
}

type GetRanklistSolutionsResponse []struct {
	Id               string             `json:"_id"`
	OrgId            string             `json:"orgId"`
	ProblemId        string             `json:"problemId"`
	ContestId        string             `json:"contestId"`
	UserId           string             `json:"userId"`
	ProblemDataHash  string             `json:"problemDataHash"`
	State            int                `json:"state"`
	SolutionDataHash string             `json:"solutionDataHash"`
	Score            float64            `json:"score"`
	Metrics          map[string]float64 `json:"metrics"`
	Status           string             `json:"status"`
	Message          string             `json:"message"`
	RunnerId         string             `json:"runnerId"`
	CreatedAt        int                `json:"createdAt"`
	SubmittedAt      int                `json:"submittedAt"`
	CompletedAt      int                `json:"completedAt"`
}

func GetRanklistSolutions(ctx context.Context, since int) (*GetRanklistSolutionsResponse, error) {
	taskId, contestId := LoadRanklistCtx(ctx)
	res := &GetRanklistSolutionsResponse{}
	raw, err := http.R().
		SetContext(ctx).
		SetQueryParam("since", fmt.Sprint(since)).
		SetResult(res).
		Get("/api/runner/ranklist/task/" + contestId + "/" + taskId + "/solutions")
	err = loadError(raw, err)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type GetRanklistParticipantsResponse []struct {
	Id        string `json:"_id" bson:"_id"`
	UserId    string `json:"userId" bson:"userId"`
	ContestId string `json:"contestId" bson:"contestId"`
	Results   map[string]struct {
		SolutionCount  int    `json:"solutionCount" bson:"solutionCount"`
		LastSolutionId string `json:"lastSolutionId" bson:"lastSolutionId"`
		LastSolution   struct {
			Score       float64 `json:"score" bson:"score"`
			Status      string  `json:"status" bson:"status"`
			CompletedAt int     `json:"completedAt" bson:"completedAt"`
		} `json:"lastSolution" bson:"lastSolution"`
	} `json:"results" bson:"results"`
	UpdatedAt int `json:"updatedAt" bson:"updatedAt"`
}

func GetRanklistParticipants(ctx context.Context, since int) (*GetRanklistParticipantsResponse, error) {
	taskId, contestId := LoadRanklistCtx(ctx)
	res := &GetRanklistParticipantsResponse{}
	raw, err := http.R().
		SetContext(ctx).
		SetQueryParam("since", fmt.Sprint(since)).
		SetResult(res).
		Get("/api/runner/ranklist/task/" + contestId + "/" + taskId + "/participants")
	err = loadError(raw, err)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type GetRanklistProblemsResponse []struct {
	Id       string   `json:"_id"`
	Title    string   `json:"title"`
	Tags     []string `json:"tags"`
	Settings struct {
		Score              float64 `json:"score"`
		Slug               string  `json:"slug"`
		SolutionCountLimit int     `json:"solutionCountLimit"`
		ShowAfter          int     `json:"showAfter"`
	} `json:"settings"`
}

func GetRanklistProblems(ctx context.Context) (*GetRanklistProblemsResponse, error) {
	taskId, contestId := LoadRanklistCtx(ctx)
	res := &GetRanklistProblemsResponse{}
	raw, err := http.R().
		SetContext(ctx).
		SetResult(res).
		Get("/api/runner/ranklist/task/" + contestId + "/" + taskId + "/problems")
	err = loadError(raw, err)
	if err != nil {
		return nil, err
	}
	return res, nil
}
