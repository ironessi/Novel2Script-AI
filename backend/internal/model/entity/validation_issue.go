package entity

import "time"

// ValidationIssue 校验问题表实体
type ValidationIssue struct {
	Id              int64     `json:"id"`
	ProjectId       int64     `json:"project_id"`
	ScriptVersionId int64     `json:"script_version_id"`
	IssueType       string    `json:"issue_type"`
	Severity        string    `json:"severity"`
	Message         string    `json:"message"`
	LocationPath    string    `json:"location_path"`
	Suggestion      string    `json:"suggestion"`
	Resolved        int       `json:"resolved"`
	CreatedAt       time.Time `json:"created_at"`
}
