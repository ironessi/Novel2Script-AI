package v1

// UploadRes 文件上传响应
type UploadRes struct {
	FileId         int64         `json:"file_id"`
	OriginalFilename string      `json:"original_filename"`
	FileSize       int64         `json:"file_size"`
	ChapterCount   int           `json:"chapter_count"`
	Chapters       []ChapterItem `json:"chapters"`
}
