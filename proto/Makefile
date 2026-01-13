all: generate-all

# 生成所有proto文件
generate-all:
	trpc create \
		-p helloworld/helloworld.proto \
		-o helloworld \
		--rpconly \
		--nogomod \
		--mock=false \

	trpc create \
		-p group_chat/group_chat.proto \
		-o group_chat \
		--rpconly \
		--nogomod \
		--mock=false \

	trpc create \
		-p conversation/conversation.proto \
		-o conversation \
		--rpconly \
		--nogomod \
		--mock=false \

	trpc create \
		-p session/session.proto \
		-o session \
		--rpconly \
		--nogomod \
		--mock=false \

	trpc create \
		-p agent/agent.proto \
		-o agent \
		--rpconly \
		--nogomod \
		--mock=false \

	trpc create \
		-p user_agent/user_agent.proto \
		-o user_agent \
		--rpconly \
		--nogomod \
		--mock=false \

	trpc create \
		-p user/user.proto \
		-o user \
		--rpconly \
		--nogomod \
		--mock=false \

	trpc create \
		-p auth/auth.proto \
		-o auth \
		--rpconly \
		--nogomod \
		--mock=false \

	trpc create \
		-p admin/admin.proto \
		-o admin \
		--rpconly \
		--nogomod \
		--mock=false \

# 清理生成的代码
clean:
	rm -rf */\*.pb.go
	rm -rf */\*.trpc.go

.PHONY: all generate-all clean