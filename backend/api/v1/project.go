package v1

// CreateProjectReq 创建项目请求
type CreateProjectReq struct {
	Title          string `json:"title"          v:"required|length:1,255#项目标题不能为空|标题长度1-255位"`
	Description    string `json:"description"`
	AdaptationMode string `json:"adaptation_mode" v:"in:screen_script,stage_play,short_video,radio_drama#改编模式不合法"`
	Visibility     string `json:"visibility"      v:"in:private,public#可见性不合法"`
}

// CreateProjectRes 创建项目响应
type CreateProjectRes struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}

// ProjectListReq 项目列表请求
type ProjectListReq struct {
	Page     int `json:"page"      d:"1"`
	PageSize int `json:"page_size" d:"10"`
}

// ProjectListRes 项目列表响应
type ProjectListRes struct {
	Total    int           `json:"total"`
	Projects []ProjectItem `json:"projects"`
}

// ProjectItem 项目列表项
type ProjectItem struct {
	Id             int64  `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	AdaptationMode string `json:"adaptation_mode"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
}

// ProjectDetailReq 项目详情请求
type ProjectDetailReq struct {
	Id int64 `json:"id" v:"required#项目ID不能为空"`
}

// ProjectDetailRes 项目详情响应
type ProjectDetailRes struct {
	Id             int64  `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	AdaptationMode string `json:"adaptation_mode"`
	Visibility     string `json:"visibility"`
	Status         string `json:"status"`
	OwnerId        int64  `json:"owner_id"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// UpdateProjectReq 更新项目请求
type UpdateProjectReq struct {
	Id             int64  `json:"id"               v:"required#项目ID不能为空"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	AdaptationMode string `json:"adaptation_mode"`
	Visibility     string `json:"visibility"`
}

// DeleteProjectReq 删除项目请求
type DeleteProjectReq struct {
	Id int64 `json:"id" v:"required#项目ID不能为空"`
}
