package rpc

import (
	"log"
	"os"

	agentpb "github.com/yzf120/elysia-backend/proto/agent"
	"trpc.group/trpc-go/trpc-go/client"
)

// AgentClient chat-agent RPC 客户端
type AgentClient struct {
	proxy agentpb.AgentServiceClientProxy
}

var defaultAgentClient *AgentClient

// InitAgentClient 初始化 chat-agent RPC 客户端
func InitAgentClient() {
	chatAgentAddr := os.Getenv("CHAT_AGENT_ADDR")
	if chatAgentAddr == "" {
		chatAgentAddr = "127.0.0.1:8081"
	}

	proxy := agentpb.NewAgentServiceClientProxy(
		client.WithTarget("ip://"+chatAgentAddr),
		client.WithTimeout(0), // 流式接口不设超时
	)

	defaultAgentClient = &AgentClient{proxy: proxy}
	log.Printf("Chat Agent RPC 客户端初始化完成，地址: %s", chatAgentAddr)
}

// GetAgentClient 获取 chat-agent RPC 客户端
func GetAgentClient() *AgentClient {
	if defaultAgentClient == nil {
		InitAgentClient()
	}
	return defaultAgentClient
}

// GetProxy 获取底层 proxy
func (c *AgentClient) GetProxy() agentpb.AgentServiceClientProxy {
	return c.proxy
}
