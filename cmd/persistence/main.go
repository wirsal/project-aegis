package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/wirsal/project-aegis/api/protos"
	"github.com/wirsal/project-aegis/internal/persistence/handler"
	"github.com/wirsal/project-aegis/internal/persistence/service"
	"github.com/wirsal/project-aegis/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("FATAL: could not load config: %v", err)
	}

	// 2. Buka listener di port gRPC yang ditentukan di config
	listener, err := net.Listen("tcp", cfg.Persistence.GRPCPort)
	if err != nil {
		log.Fatalf("FATAL: Failed to listen on port %s: %v", cfg.Persistence.GRPCPort, err)
	}

	// 3. Buat string koneksi database dari config
	dbConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)

	// 4. Buat instance dari service layer (logika bisnis)
	persistenceService, err := service.NewService(dbConnStr)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize Persistence Service: %v", err)
	}

	// 5. Buat instance dari handler gRPC dan suntikkan service ke dalamnya
	grpcHandler := handler.NewGRPCHandler(persistenceService)

	// 6. Daftarkan handler ke server gRPC
	grpcServer := grpc.NewServer()
	pb.RegisterPersistenceServer(grpcServer, grpcHandler)

	log.Printf("🚀 Persistence Service running and listening on port %s", cfg.Persistence.GRPCPort)

	// 7. Mulai melayani permintaan
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("FATAL: Failed to serve gRPC server: %v", err)
	}
}
