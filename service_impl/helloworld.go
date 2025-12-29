package service_impl

import (
	"context"
	"github.com/yzf120/elysia-backend/proto/helloworld"
)

type HelloWorldImpl struct {
}

func (h *HelloWorldImpl) Hello(ctx context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{
		Msg: "hello world",
	}, nil
}
