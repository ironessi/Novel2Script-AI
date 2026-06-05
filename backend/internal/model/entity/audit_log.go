package entity

import "time"

// AuditLog 审计日志表实体
type AuditLog struct {
	Id           int64     `json:"id"`
	UserId       int64     `json:"user_id"`
	ProjectId    int64     `json:"project_id"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceId   int64     `json:"resource_id"`
	IpAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	RequestId    string    `json:"request_id"`
	CreatedAt    time.Time `json:"created_at"`
}
