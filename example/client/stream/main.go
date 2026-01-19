package main

import (
	"context"
	"io"
	"log"

	"github.com/hoyci/gopher-chaos/example/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Falha na conexão: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	ctx := context.Background()

	stream, err := client.ListUsers(ctx, &pb.ListUserRequest{Count: 100})
	if err != nil {
		log.Fatalf("Erro ao abrir stream: %v", err)
	}

	log.Println("--- Iniciando Leitura de Stream com Caos ---")

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("Stream finalizado com sucesso. ✅")
			break
		}
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() == codes.Internal {
				log.Printf("🌪️ CAOS NO MEIO DO STREAM: %v", st.Message())
			} else {
				log.Printf("❌ Erro fatal: %v", err)
			}
			break
		}

		log.Printf("Recebido: %s", resp.User.Id)
	}
}
