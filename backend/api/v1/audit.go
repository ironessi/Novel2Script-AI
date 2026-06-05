package v1

// AuditLogListReq 审计日志列表请求
type AuditLogListReq struct {
	ProjectId int64 `json:"project_id" v:"required#项目ID不能为空"`
	Page      int   `json:"page"       d:"1"`
	PageSize  int   `json:"page_size"  d:"20"`
}

// AuditLogListRes 审计日志列表响应
type AuditLogListRes struct {
	Total int            `json:"total"`
	Logs  []AuditLogItem `json:"logs"`
}

// AuditLogItem 审计日志项
type AuditLogItem struct {
	Id           int64  `json:"id"`
	UserId       int64  `json:"user_id"`
	Action       string `json:"action"`
	ResourceType string `json:"resource_type"`
	ResourceId   int64  `json:"resource_id"`
	IpAddress    string `json:"ip_address"`
	RequestId    string `json:"request_id"`
	CreatedAt    string `json:"created_at"`
}
