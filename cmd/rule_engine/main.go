package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/wirsal/project-aegis/api/protos"
	"github.com/wirsal/project-aegis/internal/rule_engine/handler"
	"github.com/wirsal/project-aegis/internal/rule_engine/service"
	"github.com/wirsal/project-aegis/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("FATAL: could not load config: %v", err)
	}

	listener, err := net.Listen("tcp", cfg.RuleEngine.GRPCPort)
	if err != nil {
		log.Fatalf("FATAL: Failed to listen on port %s: %v", cfg.RuleEngine.GRPCPort, err)
	}

	// --- PERUBAHAN DI SINI ---

	// Buat koneksi ke Persistence Service menggunakan grpc.NewClient
	persistenceConn, err := grpc.NewClient(cfg.Persistence.GRPCPAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to Persistence Service: %v", err)
	}
	defer persistenceConn.Close()
	persistenceClient := pb.NewPersistenceClient(persistenceConn)
	log.Println("Successfully connected to Persistence Service.")

	// Buat koneksi ke Notification Service menggunakan grpc.NewClient
	notificationConn, err := grpc.NewClient(cfg.Notification.GRPCPAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to Notification Service: %v", err)
	}
	defer notificationConn.Close()
	notificationClient := pb.NewNotificationClient(notificationConn)
	log.Println("Successfully connected to Notification Service.")
	// --------------------------------------------------

	ruleEngineService, err := service.NewService("./configs/.rules.json", persistenceClient, notificationClient)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize Rule Engine Service: %v", err)
	}

	grpcHandler := handler.NewGRPCHandler(ruleEngineService)
	grpcServer := grpc.NewServer()
	pb.RegisterRuleEngineServer(grpcServer, grpcHandler)

	log.Printf("🚀 Rule Engine Service running and listening on port %s", cfg.RuleEngine.GRPCPort)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("FATAL: Failed to serve gRPC server: %v", err)
	}
}
