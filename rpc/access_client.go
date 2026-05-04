package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// AccessClient 内容安全审核 RPC 客户端（通过 HTTP 调用 elysia-access）
type AccessClient struct {
	baseURL    string
	httpClient *http.Client
}

var defaultAccessClient *AccessClient

// InitAccessClient 初始化 access 客户端
func InitAccessClient() {
	accessAddr := os.Getenv("ACCESS_SERVICE_ADDR")
	if accessAddr == "" {
		accessAddr = "http://127.0.0.1:8190"
	}
	// 确保有 http:// 前缀
	if !strings.HasPrefix(accessAddr, "http") {
		accessAddr = "http://" + accessAddr
	}

	defaultAccessClient = &AccessClient{
		baseURL: accessAddr,
		httpClient: &http.Client{
			Timeout: 35 * time.Second,
		},
	}
	log.Printf("Access 内容安全审核客户端初始化完成，地址: %s", accessAddr)
}

// GetAccessClient 获取 access 客户端
func GetAccessClient() *AccessClient {
	if defaultAccessClient == nil {
		InitAccessClient()
	}
	return defaultAccessClient
}

// TextModerationRequest 文本审核请求
type TextModerationRequest struct {
	Content string `json:"content"`
	BizType string `json:"biz_type,omitempty"`
	UserId  string `json:"user_id,omitempty"`
	DataId  string `json:"data_id,omitempty"`
}

// TextModerationResponse 文本审核响应
type TextModerationResponse struct {
	Code       int32  `json:"code"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
	Label      string `json:"label"`
	Score      int32  `json:"score"`
	RequestId  string `json:"request_id"`
}

// ImageModerationRequest 图片审核请求
type ImageModerationRequest struct {
	FileUrl     string `json:"file_url,omitempty"`
	FileContent []byte `json:"file_content,omitempty"`
	BizType     string `json:"biz_type,omitempty"`
	UserId      string `json:"user_id,omitempty"`
	DataId      string `json:"data_id,omitempty"`
}

// ImageModerationResponse 图片审核响应
type ImageModerationResponse struct {
	Code       int32  `json:"code"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
	Label      string `json:"label"`
	Score      int32  `json:"score"`
	RequestId  string `json:"request_id"`
}

// TextModeration 调用文本审核接口
func (c *AccessClient) TextModeration(ctx context.Context, req *TextModerationRequest) (*TextModerationResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/content-security/text", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("调用文本审核接口失败: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("文本审核接口返回非200: %d", httpResp.StatusCode)
	}

	var resp TextModerationResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("解析文本审核响应失败: %w", err)
	}

	return &resp, nil
}

// ImageModeration 调用图片审核接口
func (c *AccessClient) ImageModeration(ctx context.Context, req *ImageModerationRequest) (*ImageModerationResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/content-security/image", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("调用图片审核接口失败: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("图片审核接口返回非200: %d", httpResp.StatusCode)
	}

	var resp ImageModerationResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("解析图片审核响应失败: %w", err)
	}

	return &resp, nil
}
