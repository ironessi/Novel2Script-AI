package entity

import "time"

// NovelProject 小说项目表实体
type NovelProject struct {
	Id             int64      `json:"id"`
	OwnerId        int64      `json:"owner_id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	AdaptationMode string     `json:"adaptation_mode"`
	Visibility     string     `json:"visibility"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}
