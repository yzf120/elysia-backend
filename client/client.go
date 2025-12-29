package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient 封装HTTP客户端功能
type HTTPClient struct {
	client  *http.Client
	baseURL string
}

// NewHTTPClient 创建新的HTTP客户端
func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

// Get 发送GET请求
func (c *HTTPClient) Get(path string, headers map[string]string) ([]byte, error) {
	return c.doRequest("GET", path, nil, headers)
}

// Post 发送POST请求
func (c *HTTPClient) Post(path string, body interface{}, headers map[string]string) ([]byte, error) {
	return c.doRequest("POST", path, body, headers)
}

// doRequest 执行HTTP请求
func (c *HTTPClient) doRequest(method, path string, body interface{}, headers map[string]string) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置默认头
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	if resp.StatusCode >= 400 {
		return respBody, fmt.Errorf("HTTP错误: %d - %s", resp.StatusCode, resp.Status)
	}

	return respBody, nil
}
