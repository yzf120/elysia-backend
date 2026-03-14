package rpc

import (
	"context"
	"log"
	"os"
	"time"

	agent_session "github.com/yzf120/elysia-session/proto/agent_session"
	conversation "github.com/yzf120/elysia-session/proto/conversation"
	"trpc.group/trpc-go/trpc-go/client"
)

// SessionClient elysia-session 服务 tRPC 客户端
type SessionClient struct {
	agentSessionProxy agent_session.AgentSessionServiceClientProxy
	conversationProxy conversation.ConversationServiceClientProxy
}

var defaultSessionClient *SessionClient

// InitSessionClient 初始化 session 服务 tRPC 客户端
func InitSessionClient() {
	sessionAgentAddr := os.Getenv("SESSION_AGENT_ADDR")
	if sessionAgentAddr == "" {
		sessionAgentAddr = "127.0.0.1:8002"
	}
	sessionConvAddr := os.Getenv("SESSION_CONV_ADDR")
	if sessionConvAddr == "" {
		sessionConvAddr = "127.0.0.1:8003"
	}

	defaultSessionClient = &SessionClient{
		agentSessionProxy: agent_session.NewAgentSessionServiceClientProxy(
			client.WithTarget("ip://"+sessionAgentAddr),
			client.WithTimeout(10000*time.Millisecond),
		),
		conversationProxy: conversation.NewConversationServiceClientProxy(
			client.WithTarget("ip://"+sessionConvAddr),
			client.WithTimeout(10000*time.Millisecond),
		),
	}
	log.Printf("Session tRPC 客户端初始化完成，AgentSession: %s, Conversation: %s", sessionAgentAddr, sessionConvAddr)
}

// GetSessionClient 获取 session 服务客户端
func GetSessionClient() *SessionClient {
	if defaultSessionClient == nil {
		InitSessionClient()
	}
	return defaultSessionClient
}

// CreateSession 创建会话
func (c *SessionClient) CreateSession(ctx context.Context, req *agent_session.CreateSessionRequest) (*agent_session.CreateSessionResponse, error) {
	return c.agentSessionProxy.CreateSession(ctx, req)
}

// CreateConversation 创建对话消息
func (c *SessionClient) CreateConversation(ctx context.Context, req *conversation.CreateConversationRequest) (*conversation.CreateConversationResponse, error) {
	return c.conversationProxy.CreateConversation(ctx, req)
}

// ListSessionsByUser 获取用户会话列表
func (c *SessionClient) ListSessionsByUser(ctx context.Context, req *agent_session.ListSessionsByUserRequest) (*agent_session.ListSessionsByUserResponse, error) {
	return c.agentSessionProxy.ListSessionsByUser(ctx, req)
}

// ListConversations 获取会话消息列表
func (c *SessionClient) ListConversations(ctx context.Context, req *conversation.ListConversationsRequest) (*conversation.ListConversationsResponse, error) {
	return c.conversationProxy.ListConversations(ctx, req)
}
