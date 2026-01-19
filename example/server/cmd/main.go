package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/hoyci/gopher-chaos/example/pb"
	"github.com/hoyci/gopher-chaos/example/server/internal/handlers"
	"github.com/hoyci/gopher-chaos/example/server/internal/repositories"
	"github.com/hoyci/gopher-chaos/example/server/internal/services"
	"github.com/hoyci/gopher-chaos/pkg/chaos"
	"github.com/hoyci/gopher-chaos/pkg/interceptors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func main() {
	repo := repositories.NewMemoryUserRepository()
	useCase := &services.UserUseCase{Repo: repo}
	grpcHandler := &handlers.UserGRPCHandler{UseCase: useCase}
	logger := log.New(os.Stdout, "[SERVER] ", log.LstdFlags|log.Lmicroseconds)

	chaosCfg := chaos.ChaosConfig{
		Probability: 0.1,
		Latency: chaos.ChaosConfigLatency{
			Min: 1 * time.Millisecond,
			Max: 1 * time.Second,
		},
		Error: codes.Internal,
	}

	chaosEngine := chaos.NewChaos(chaosCfg, chaos.WithLogger(logger))
	chaosInterceptor := interceptors.NewInterceptor(chaosEngine)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(chaosInterceptor.UnaryInterceptor),
		grpc.StreamInterceptor(chaosInterceptor.StreamInterceptor),
	)

	pb.RegisterUserServiceServer(s, grpcHandler)

	log.Println("Server gRPC running with Chaos Engine on port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
