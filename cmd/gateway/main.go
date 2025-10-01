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

func main() {
	// Load configuration from file
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("FATAL: could not load config: %v", err)
	}

	// Use the value from config, not a constant
	// PERUBAHAN: Gunakan grpc.NewClient
	conn, err := grpc.NewClient(cfg.Gateway.RuleEngineAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	ruleEngineClient := pb.NewRuleEngineClient(conn)
	log.Println("Successfully connected to Rule Engine Service.")

	gatewayService := service.NewGatewayService(ruleEngineClient)
	tcpHandler := handler.NewTCPHandler(gatewayService)

	// Use the port from config
	tcpHandler.StartServer(cfg.Gateway.TCPPort)
}
