package client

import "context"

type RegisterRequest struct {
	Name              string   `json:"name"`
	Labels            []string `json:"labels"`
	Version           string   `json:"version"`
	RegistrationToken string   `json:"registrationToken"`
}

type RegisterResponse struct {
	RunnerId  string `json:"runnerId"`
	RunnerKey string `json:"runnerKey"`
}

func Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	res := &RegisterResponse{}
	raw, err := http.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(res).
		Post("/api/runner/register")
	err = loadError(raw, err)
	if err != nil {
		return nil, err
	}
	return res, nil
}
