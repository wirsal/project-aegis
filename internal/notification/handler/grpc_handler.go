package handler

import (
	"context"
	"log"

	pb "github.com/wirsal/project-aegis/api/protos"
)

type NotificationService interface {
	TriggerRiskAlert(ctx context.Context, req *pb.RiskAlertRequest) error
}

type GRPCHandler struct {
	pb.UnimplementedNotificationServer
	service NotificationService
}

func NewGRPCHandler(svc NotificationService) *GRPCHandler {
	return &GRPCHandler{
		service: svc,
	}
}

func (h *GRPCHandler) TriggerRiskAlert(ctx context.Context, in *pb.RiskAlertRequest) (*pb.NotificationAck, error) {
	log.Printf("gRPC handler received TriggerRiskAlert for TrxKey: %s", in.GetTransactionData().GetTrxKey())

	// 3. Panggil metode service yang sudah diperbarui
	err := h.service.TriggerRiskAlert(ctx, in)
	if err != nil {
		log.Printf("ERROR from service layer: %v", err)
		return &pb.NotificationAck{Success: false, DeliveryStatus: "FAILED"}, nil
	}

	return &pb.NotificationAck{Success: true, DeliveryStatus: "SENT"}, nil
}
