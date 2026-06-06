package entity

import "time"

// NovelSourceFile 原始文件表实体
type NovelSourceFile struct {
	Id               int64     `json:"id"`
	ProjectId        int64     `json:"project_id"`
	OwnerId          int64     `json:"owner_id"`
	OriginalFilename string    `json:"original_filename"`
	StoragePath      string    `json:"storage_path"`
	FileHash         string    `json:"file_hash"`
	FileSize         int64     `json:"file_size"`
	MimeType         string    `json:"mime_type"`
	ScanStatus       string    `json:"scan_status"`
	CreatedAt        time.Time `json:"created_at"`
}
