package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/wirsal/project-aegis/api/protos"
	"github.com/wirsal/project-aegis/internal/gateway/handler"
	"github.com/wirsal/project-aegis/internal/gateway/service"
	"github.com/wirsal/project-aegis/pkg/config"
)

// cmd/gateway/main.go

func main() {
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("FATAL: could not load config: %v", err)
	}

	// Buat koneksi ke Rule Engine Service
	ruleEngineConn, err := grpc.NewClient(cfg.Gateway.RuleEngineAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to Rule Engine gRPC server: %v", err)
	}
	defer ruleEngineConn.Close()
	ruleEngineClient := pb.NewRuleEngineClient(ruleEngineConn)
	log.Println("Successfully connected to Rule Engine Service.")

	// --- TAMBAHAN: Buat koneksi ke Persistence Service ---
	persistenceConn, err := grpc.NewClient(cfg.Persistence.GRPCPAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to Persistence Service: %v", err)
	}
	defer persistenceConn.Close()
	persistenceClient := pb.NewPersistenceClient(persistenceConn)
	log.Println("Successfully connected to Persistence Service.")
	// --------------------------------------------------

	// Suntikkan KEDUA klien ke dalam Gateway Service
	gatewayService := service.NewGatewayService(ruleEngineClient, persistenceClient)
	tcpHandler := handler.NewTCPHandler(gatewayService)

	tcpHandler.StartServer(cfg.Gateway.TCPPort)
}
