package client

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

type APIError struct {
	Message    string `json:"message"`
	ErrorName  string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

func (e *APIError) Error() string {
	return e.Message
}

func loadError(res *resty.Response, err error) error {
	if err != nil {
		return err
	}
	if res.IsError() {
		apiError := &APIError{}
		err := json.Unmarshal(res.Body(), apiError)
		if err != nil {
			return err
		}
		return apiError
	}
	return nil
}
