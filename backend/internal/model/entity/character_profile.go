package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// CharacterProfile 人物档案表实体
type CharacterProfile struct {
	Id            int64           `json:"id"`
	ProjectId     int64           `json:"project_id"`
	CharacterKey  string          `json:"character_key"`
	Name          string          `json:"name"`
	Aliases       JSON            `json:"aliases"`
	RoleType      string          `json:"role_type"`
	Description   string          `json:"description"`
	Personality   JSON            `json:"personality"`
	Relationships JSON            `json:"relationships"`
	SourceRefs    JSON            `json:"source_refs"`
	Confidence    float64         `json:"confidence"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// JSON 自定义 JSON 类型，支持数据库 JSON 字段
type JSON json.RawMessage

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "[]", nil
	}
	return string(j), nil
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = JSON("[]")
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		*j = JSON("[]")
		return nil
	}
	*j = JSON(bytes)
	return nil
}

func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("[]"), nil
	}
	return []byte(j), nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	*j = JSON(data)
	return nil
}
