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

// TaskResponse 任务响应
type TaskResponse struct {
	TaskId string `json:"task_id"`
	Status string `json:"status"`
}

// Analyze 发送分析请求到 AI 服务
func (c *AIServiceClient) Analyze(ctx context.Context, req *AnalyzeRequest) (*TaskResponse, error) {
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

	var result TaskResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
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
