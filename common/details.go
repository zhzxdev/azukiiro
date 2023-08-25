package common

type SolutionDetailsTest struct {
	Name    string  `json:"name"`
	Score   float64 `json:"score"`
	Status  string  `json:"status"`
	Summary string  `json:"summary"`
}

type SolutionDetailsJob struct {
	Name       string                 `json:"name"`
	Score      float64                `json:"score"`
	ScoreScale float64                `json:"scoreScale"`
	Status     string                 `json:"status"`
	Tests      []*SolutionDetailsTest `json:"tests"`
	Summary    string                 `json:"summary"`
}

type SolutionDetails struct {
	Jobs    []*SolutionDetailsJob `json:"jobs"`
	Summary string                `json:"summary"`
}
