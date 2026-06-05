package entity

import "time"

// NovelChapter 小说章节表实体
type NovelChapter struct {
	Id           int64     `json:"id"`
	ProjectId    int64     `json:"project_id"`
	ChapterIndex int       `json:"chapter_index"`
	ChapterTitle string    `json:"chapter_title"`
	Content      string    `json:"content"`
	ContentHash  string    `json:"content_hash"`
	CreatedAt    time.Time `json:"created_at"`
}
