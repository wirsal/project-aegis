package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/wirsal/project-aegis/api/protos"
	"github.com/wirsal/project-aegis/internal/notification/handler"
	"github.com/wirsal/project-aegis/internal/notification/service"
	"github.com/wirsal/project-aegis/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("FATAL: could not load config: %v", err)
	}

	listener, err := net.Listen("tcp", cfg.Notification.GRPCPort)
	if err != nil {
		log.Fatalf("FATAL: Failed to listen on port %s: %v", cfg.Notification.GRPCPort, err)
	}

	// Buat instance service dan berikan URL API eksternal dari config
	notificationService := service.NewService(cfg.Notification.ExternalAPIURL)
	grpcHandler := handler.NewGRPCHandler(notificationService)

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServer(grpcServer, grpcHandler)

	log.Printf("🚀 Notification Service running and listening on port %s", cfg.Notification.GRPCPort)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("FATAL: Failed to serve gRPC server: %v", err)
	}
}
