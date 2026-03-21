package rpc

import (
	"log"
	"os"
	"time"

	llmpb "github.com/yzf120/elysia-llm-tool/proto/llm"
	"trpc.group/trpc-go/trpc-go/client"
)

// LLMToolClient llm-tool RPC 客户端
type LLMToolClient struct {
	proxy llmpb.LLMServiceClientProxy
}

var defaultLLMToolClient *LLMToolClient

// InitLLMToolClient 初始化 llm-tool RPC 客户端
func InitLLMToolClient() {
	llmToolAddr := os.Getenv("LLM_TOOL_ADDR")
	if llmToolAddr == "" {
		llmToolAddr = "127.0.0.1:9001"
	}

	proxy := llmpb.NewLLMServiceClientProxy(
		client.WithTarget("ip://"+llmToolAddr),
		client.WithTimeout(30*time.Second),
	)

	defaultLLMToolClient = &LLMToolClient{proxy: proxy}
	log.Printf("LLM Tool RPC 客户端初始化完成，地址: %s", llmToolAddr)
}

// GetLLMToolClient 获取 llm-tool RPC 客户端
func GetLLMToolClient() *LLMToolClient {
	if defaultLLMToolClient == nil {
		InitLLMToolClient()
	}
	return defaultLLMToolClient
}

// GetProxy 获取底层 proxy
func (c *LLMToolClient) GetProxy() llmpb.LLMServiceClientProxy {
	return c.proxy
}
