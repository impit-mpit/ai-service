package router

import (
	"context"
	"log"
	"net"
	aiv1 "neuro-most/ai-service/gen/go/ai/v1"
	"neuro-most/ai-service/internal/adapters/action"
	"neuro-most/ai-service/internal/adapters/vllm"
	"neuro-most/ai-service/internal/usecase"

	"google.golang.org/grpc"
)

type Router struct {
	vllm vllm.Vllm
	aiv1.UnimplementedAIServiceServer
}

func NewRouter(vllm vllm.Vllm) Router {
	return Router{vllm: vllm}
}

func (r *Router) Listen() {
	port := ":3001"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts = []grpc.ServerOption{}
	srv := grpc.NewServer(opts...)
	aiv1.RegisterAIServiceServer(srv, r)

	log.Printf("Starting gRPC server on port %s\n", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (r *Router) Chat(ctx context.Context, req *aiv1.CreateChatRequest) (*aiv1.ChatResponse, error) {
	var (
		uc  = usecase.NewChatUseCase(r.vllm)
		act = action.NewChatAction(uc)
	)

	return act.Execute(ctx, req)
}
