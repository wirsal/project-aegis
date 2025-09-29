package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/wirsal/project-aegis/api/protos"
	"github.com/wirsal/project-aegis/internal/rule_engine/handler"
	"github.com/wirsal/project-aegis/internal/rule_engine/service"
	"github.com/wirsal/project-aegis/pkg/config"
)

func main() {
	// 2. Muat konfigurasi dari file
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("FATAL: could not load config: %v", err)
	}

	// 3. Gunakan port dari config
	listener, err := net.Listen("tcp", cfg.RuleEngine.GRPCPort)
	if err != nil {
		log.Fatalf("FATAL: Failed to listen on port %s: %v", cfg.RuleEngine.GRPCPort, err)
	}

	ruleEngineService, err := service.NewService("./configs/rules.json")
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
