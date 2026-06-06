package entity

import "time"

// AiTask AI 任务表实体
type AiTask struct {
	Id           int64      `json:"id"`
	ProjectId    int64      `json:"project_id"`
	OwnerId      int64      `json:"owner_id"`
	TaskType     string     `json:"task_type"`
	Status       string     `json:"status"`
	Progress     int        `json:"progress"`
	CurrentStep  string     `json:"current_step"`
	ErrorMessage string     `json:"error_message"`
	RetryCount   int        `json:"retry_count"`
	StartedAt    *time.Time `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
