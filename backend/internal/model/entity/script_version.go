package entity

import "time"

// ScriptVersion 剧本版本表实体
type ScriptVersion struct {
	Id                int64     `json:"id"`
	ProjectId         int64     `json:"project_id"`
	OwnerId           int64     `json:"owner_id"`
	VersionNo         int       `json:"version_no"`
	YamlContent       string    `json:"yaml_content"`
	YamlHash          string    `json:"yaml_hash"`
	ValidationStatus  string    `json:"validation_status"`
	HallucinationRisk string    `json:"hallucination_risk"`
	SafetyRisk        string    `json:"safety_risk"`
	CreatedBy         string    `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
}
