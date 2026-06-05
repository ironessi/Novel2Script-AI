package entity

import "time"

// PlotEvent 剧情事件表实体
type PlotEvent struct {
	Id          int64     `json:"id"`
	ProjectId   int64     `json:"project_id"`
	EventKey    string    `json:"event_key"`
	ChapterIndex int      `json:"chapter_index"`
	TriggerText string    `json:"trigger_text"`
	ActionText  string    `json:"action_text"`
	ResultText  string    `json:"result_text"`
	Importance  string    `json:"importance"`
	SourceRefs  JSON      `json:"source_refs"`
	Confidence  float64   `json:"confidence"`
	CreatedAt   time.Time `json:"created_at"`
}
