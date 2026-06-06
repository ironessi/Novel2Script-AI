package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"novel2script-backend/internal/config"
)

// AIServiceClient AI 服务客户端
type AIServiceClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewAIServiceClient 创建 AI 服务客户端
func NewAIServiceClient() *AIServiceClient {
	return &AIServiceClient{
		baseURL: config.C.AI.URL,
		token:   config.C.AI.Token,
		httpClient: &http.Client{
			Timeout: 300 * time.Second,
		},
	}
}

// AnalyzeRequest 分析请求
type AnalyzeRequest struct {
	ProjectId      int64          `json:"project_id"`
	Chapters       []ChapterInput `json:"chapters"`
	AdaptationMode string         `json:"adaptation_mode"`
}

// ChapterInput 章节输入
type ChapterInput struct {
	ChapterIndex int    `json:"chapter_index"`
	Title        string `json:"title"`
	Content      string `json:"content"`
}

// AnalyzeResponse 完整分析响应
type AnalyzeResponse struct {
	Status       string           `json:"status"`
	YamlContent  string           `json:"yaml_content"`
	Characters   []CharacterData  `json:"characters"`
	PlotEvents   []PlotEventData  `json:"plot_events"`
	Scenes       []SceneData      `json:"scenes"`
	Validation   ValidationResult `json:"validation"`
	Hallucination RiskResult      `json:"hallucination"`
	Safety       SafetyResult     `json:"safety"`
}

// CharacterData 人物数据
type CharacterData struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Aliases     []string       `json:"aliases"`
	Role        string         `json:"role"`
	Description string         `json:"description"`
	Personality []string       `json:"personality"`
	Relationships []Relationship `json:"relationships"`
	SourceTrace []SourceTrace  `json:"source_trace"`
	Confidence  float64        `json:"confidence"`
}

// Relationship 人物关系
type Relationship struct {
	Target   string `json:"target"`
	Relation string `json:"relation"`
}

// SourceTrace 原文溯源
type SourceTrace struct {
	ChapterIndex   int `json:"chapter_index"`
	ParagraphStart int `json:"paragraph_start"`
	ParagraphEnd   int `json:"paragraph_end"`
}

// PlotEventData 剧情事件数据
type PlotEventData struct {
	Id                string        `json:"id"`
	EventKey          string        `json:"event_key"`
	ChapterIndex      int           `json:"chapter_index"`
	Trigger           string        `json:"trigger"`
	Action            string        `json:"action"`
	Result            string        `json:"result"`
	Importance        string        `json:"importance"`
	CharactersInvolved []string     `json:"characters_involved"`
	SourceTrace       SourceTrace   `json:"source_trace"`
	Confidence        float64       `json:"confidence"`
}

// SceneData 场景数据
type SceneData struct {
	Id              string       `json:"id"`
	Title           string       `json:"title"`
	Order           int          `json:"order"`
	Time            string       `json:"time"`
	Location        string       `json:"location"`
	Characters      []string     `json:"characters"`
	Summary         string       `json:"summary"`
	SourceTrace     []SourceTrace `json:"source_trace"`
	RelatedEvents   []string     `json:"related_events"`
	Actions         []ActionData `json:"actions"`
	Dialogues       []DialogueData `json:"dialogues"`
	AdaptationNotes []AdaptationNote `json:"adaptation_notes"`
	Confidence      float64      `json:"confidence"`
}

// ActionData 动作数据
type ActionData struct {
	Character   string `json:"character"`
	Description string `json:"description"`
}

// DialogueData 对白数据
type DialogueData struct {
	Character string `json:"character"`
	Line      string `json:"line"`
	Emotion   string `json:"emotion"`
}

// AdaptationNote 改编说明
type AdaptationNote struct {
	Source string `json:"source"`
	Change string `json:"change"`
	Reason string `json:"reason"`
}

// ValidationResult 校验结果
type ValidationResult struct {
	Valid  bool            `json:"valid"`
	Issues []ValidationIssue `json:"issues"`
}

// ValidationIssue 校验问题
type ValidationIssue struct {
	Type         string `json:"type"`
	Severity     string `json:"severity"`
	Message      string `json:"message"`
	LocationPath string `json:"location_path"`
	Suggestion   string `json:"suggestion"`
}

// RiskResult 风险结果
type RiskResult struct {
	RiskLevel string      `json:"risk_level"`
	Issues    []RiskIssue `json:"issues"`
}

// RiskIssue 风险问题
type RiskIssue struct {
	Type       string `json:"type"`
	Severity   string `json:"severity"`
	Location   string `json:"location"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
}

// SafetyResult 安全审查结果
type SafetyResult struct {
	Safe      bool        `json:"safe"`
	RiskLevel string      `json:"risk_level"`
	Issues    []RiskIssue `json:"issues"`
}

// Analyze 发送完整分析请求到 AI 服务（同步等待结果）
func (c *AIServiceClient) Analyze(ctx context.Context, req *AnalyzeRequest) (*AnalyzeResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/ai/analyze", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("AI service request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result AnalyzeResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return &result, nil
}

// HealthCheck 健康检查
func (c *AIServiceClient) HealthCheck(ctx context.Context) error {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("AI service health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("AI service health check returned status %d", resp.StatusCode)
	}

	return nil
}

// doPost 通用 POST 请求
func (c *AIServiceClient) doPost(ctx context.Context, path string, req interface{}, result interface{}) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("AI service request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("AI service returned %d: %s", resp.StatusCode, string(respBody))
	}
	return json.Unmarshal(respBody, result)
}

// ValidateYAMLRequest 校验请求
type ValidateYAMLRequest struct {
	YamlContent string `json:"yaml_content"`
}

// ValidateYAMLResponse 校验响应
type ValidateYAMLResponse struct {
	Valid  bool              `json:"valid"`
	Issues []ValidationIssue `json:"issues"`
}

// ValidateYAML 调用 AI 服务校验 YAML
func (c *AIServiceClient) ValidateYAML(ctx context.Context, yamlContent string) (*ValidateYAMLResponse, error) {
	var result ValidateYAMLResponse
	err := c.doPost(ctx, "/ai/validate-yaml", ValidateYAMLRequest{YamlContent: yamlContent}, &result)
	return &result, err
}

// HallucinationRequest 幻觉检测请求
type HallucinationRequest struct {
	YamlContent string          `json:"yaml_content"`
	Chapters    []ChapterInput  `json:"chapters"`
	Characters  []CharacterData `json:"characters"`
}

// HallucinationResponse 幻觉检测响应
type HallucinationResponse = RiskResult

// CheckHallucination 调用 AI 服务检测幻觉
func (c *AIServiceClient) CheckHallucination(ctx context.Context, yamlContent string, chapters []ChapterInput, characters []CharacterData) (*HallucinationResponse, error) {
	var result HallucinationResponse
	err := c.doPost(ctx, "/ai/check-hallucination", HallucinationRequest{
		YamlContent: yamlContent,
		Chapters:    chapters,
		Characters:  characters,
	}, &result)
	return &result, err
}

// SafetyRequest 安全审查请求
type SafetyRequest struct {
	YamlContent string `json:"yaml_content"`
}

// SafetyResponse 安全审查响应
type SafetyResponse = SafetyResult

// CheckSafety 调用 AI 服务安全审查
func (c *AIServiceClient) CheckSafety(ctx context.Context, yamlContent string) (*SafetyResponse, error) {
	var result SafetyResponse
	err := c.doPost(ctx, "/ai/check-safety", SafetyRequest{YamlContent: yamlContent}, &result)
	return &result, err
}

// RepairYAMLRequest 修复请求
type RepairYAMLRequest struct {
	YamlContent string            `json:"yaml_content"`
	Issues      []ValidationIssue `json:"issues"`
}

// RepairYAMLResponse 修复响应
type RepairYAMLResponse struct {
	FixedYaml string            `json:"fixed_yaml"`
	Changes   []RepairChange    `json:"changes"`
	Success   bool              `json:"success"`
}

// RepairChange 修复变更
type RepairChange struct {
	Location string `json:"location"`
	Original string `json:"original"`
	Fixed    string `json:"fixed"`
	Reason   string `json:"reason"`
}

// RepairYAML 调用 AI 服务修复 YAML
func (c *AIServiceClient) RepairYAML(ctx context.Context, yamlContent string, issues []ValidationIssue) (*RepairYAMLResponse, error) {
	var result RepairYAMLResponse
	err := c.doPost(ctx, "/ai/repair-yaml", RepairYAMLRequest{
		YamlContent: yamlContent,
		Issues:      issues,
	}, &result)
	return &result, err
}
