package handler

import (
	"context"
	"log"

	pb "github.com/wirsal/project-aegis/api/protos"
	"github.com/wirsal/project-aegis/internal/notification/service"
)

type GRPCHandler struct {
	pb.UnimplementedNotificationServer
	service *service.Service
}

func NewGRPCHandler(svc *service.Service) *GRPCHandler {
	return &GRPCHandler{
		service: svc,
	}
}

func (h *GRPCHandler) SendRiskNotification(ctx context.Context, in *pb.RiskResult) (*pb.NotificationAck, error) {
	log.Printf("gRPC handler received SendRiskNotification for TrxKey: %s", in.TrxKey)
	err := h.service.SendRiskNotification(ctx, in)
	if err != nil {
		return &pb.NotificationAck{Success: false, DeliveryStatus: "FAILED"}, nil
	}
	return &pb.NotificationAck{Success: true, DeliveryStatus: "SENT"}, nil
}
