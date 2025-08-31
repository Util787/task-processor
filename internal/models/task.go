package models

type Task struct {
	Id         string `json:"id"`
	Payload    string `json:"payload"`
	MaxRetries int    `json:"max_retries"`
}

type TaskState string

const (
	StateQueued  TaskState = "queued"
	StateRunning TaskState = "running"
	StateDone    TaskState = "done"
	StateFailed  TaskState = "failed"
)
