package v1

// CreateTaskReq 创建任务请求
type CreateTaskReq struct {
	ProjectId int64  `json:"project_id"`
	TaskType  string `json:"task_type" v:"required|in:full_generate,character_extract,plot_extract,scene_split,validate#任务类型不合法"`
}

// CreateTaskRes 创建任务响应
type CreateTaskRes struct {
	TaskId int64  `json:"task_id"`
	Status string `json:"status"`
}

// TaskStatusReq 任务状态请求
type TaskStatusReq struct {
	TaskId int64 `json:"task_id" v:"required#任务ID不能为空"`
}

// TaskStatusRes 任务状态响应
type TaskStatusRes struct {
	Id           int64  `json:"id"`
	ProjectId    int64  `json:"project_id"`
	TaskType     string `json:"task_type"`
	Status       string `json:"status"`
	Progress     int    `json:"progress"`
	CurrentStep  string `json:"current_step"`
	ErrorMessage string `json:"error_message"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
