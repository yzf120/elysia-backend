package router

import "github.com/gorilla/mux"

func Init() {
}
func RegisterRouter(router *mux.Router) {
	// 统一鉴权
	middlewares := []mux.MiddlewareFunc{}
	router.Use(middlewares...)
	registerApiRouter(router)
}
func registerApiRouter(router *mux.Router) {
	// 注册认证相关接口
	registerAuth(router)
	// 注册会话相关的
	registerConversation(router)
	// 注册群聊相关的
	registerGroupChat(router)
	// 注册会话相关的
	registerSession(router)
	// 注册用户智能体相关接口
	registerUserAgent(router)
	// 注册智能体相关
	registerAgent(router)
	// 注册用户相关接口
	registerUser(router)
}
