package types

type Job struct {
	Code      string `json:"code"`
	Language  string `json:"language"`
	JobID     string `json:"job_id"`
	ProblemId uint   `json:"problem_id"`
}
