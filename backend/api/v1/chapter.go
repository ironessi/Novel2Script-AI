package v1

// ChapterListReq 章节列表请求
type ChapterListReq struct {
	ProjectId int64 `json:"project_id" v:"required#项目ID不能为空"`
}

// ChapterListRes 章节列表响应
type ChapterListRes struct {
	Chapters []ChapterItem `json:"chapters"`
}

// ChapterItem 章节列表项
type ChapterItem struct {
	Id           int64  `json:"id"`
	ChapterIndex int    `json:"chapter_index"`
	ChapterTitle string `json:"chapter_title"`
	ContentHash  string `json:"content_hash"`
	CreatedAt    string `json:"created_at"`
}

// ChapterDetailReq 章节详情请求
type ChapterDetailReq struct {
	ProjectId int64 `json:"project_id" v:"required#项目ID不能为空"`
	ChapterId int64 `json:"chapter_id" v:"required#章节ID不能为空"`
}

// ChapterDetailRes 章节详情响应
type ChapterDetailRes struct {
	Id           int64  `json:"id"`
	ChapterIndex int    `json:"chapter_index"`
	ChapterTitle string `json:"chapter_title"`
	Content      string `json:"content"`
}
