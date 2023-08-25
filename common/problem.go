package common

import "encoding/json"

type ProblemConfigJudge struct {
	Adapter string          `json:"adapter"`
	Config  json.RawMessage `json:"config"`
}

type ProblemConfig struct {
	Label    string             `json:"label"`
	Solution interface{}        `json:"solution"`
	Judge    ProblemConfigJudge `json:"judge"`
	Submit   interface{}        `json:"submit"`
}
