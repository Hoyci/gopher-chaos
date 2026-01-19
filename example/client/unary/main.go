package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hoyci/gopher-chaos/example/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	workers       = 10 // Quantas goroutines em paralelo
	reqsPerWorker = 50 // Quantas requisições cada worker fará
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Falha na conexão: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	var (
		successCount uint64
		chaosCount   uint64
		totalTime    int64 // em milissegundos
	)

	var wg sync.WaitGroup
	start := time.Now()

	log.Printf("Iniciando teste: %d workers, %d requisições cada\n", workers, reqsPerWorker)

	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < reqsPerWorker; i++ {
				reqStart := time.Now()

				// Cada request precisa de seu próprio contexto com timeout
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

				_, err := client.CreateUser(ctx, &pb.CreateUserRequest{
					Name:  fmt.Sprintf("User-%d-%d", workerID, i),
					Email: fmt.Sprintf("test-%d@chaos.com", workerID*i),
					Age:   20,
				})

				if err != nil {
					// Verificamos se o erro foi injetado pelo nosso Chaos Engine
					st, ok := status.FromError(err)
					if ok && (st.Code() == codes.Internal || st.Code() == codes.Unavailable) {
						atomic.AddUint64(&chaosCount, 1)
						log.Printf("[Worker %d] CAOS DETECTADO: %v", workerID, st.Message())
					} else {
						log.Printf("[Worker %d] Erro inesperado: %v", workerID, err)
					}
				} else {
					atomic.AddUint64(&successCount, 1)
				}

				atomic.AddInt64(&totalTime, int64(time.Since(reqStart).Milliseconds()))
				cancel()
			}
		}(w)
	}

	wg.Wait()
	duration := time.Since(start)

	// --- RELATÓRIO FINAL ---
	printReport(successCount, chaosCount, duration, totalTime)
}

func printReport(success, chaos uint64, duration time.Duration, totalTime int64) {
	total := success + chaos
	avgLat := float64(totalTime) / float64(total)

	fmt.Println("\n" + "---" + " RESULTADO DO TESTE DE CARGA " + "---")
	fmt.Printf("Tempo Total:       %v\n", duration)
	fmt.Printf("Requisições:       %d\n", total)
	fmt.Printf("Sucessos:          %d ✅\n", success)
	fmt.Printf("Falhas (Caos):     %d 🌪️\n", chaos)
	fmt.Printf("Latência Média:    %.2fms\n", avgLat)
	fmt.Printf("Sucesso Rate:      %.2f%%\n", (float64(success)/float64(total))*100)
	fmt.Println("---------------------------------------")
}
