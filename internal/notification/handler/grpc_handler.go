package handler

import (
	"context"
	"log"

	pb "github.com/wirsal/project-aegis/api/protos"
)

// 1. Buat interface yang mendefinisikan apa yang dibutuhkan handler dari service.
type NotificationService interface {
	SendRiskNotification(ctx context.Context, riskData *pb.RiskResult) error
}

// GRPCHandler sekarang bergantung pada interface, bukan struct.
type GRPCHandler struct {
	pb.UnimplementedNotificationServer
	service NotificationService // <-- Bergantung pada interface
}

// NewGRPCHandler menerima interface sebagai argumen.
func NewGRPCHandler(svc NotificationService) *GRPCHandler {
	return &GRPCHandler{
		service: svc,
	}
}

func (h *GRPCHandler) SendRiskNotification(ctx context.Context, in *pb.RiskResult) (*pb.NotificationAck, error) {
	log.Printf("gRPC handler received SendRiskNotification for RRN: %s", in.TrxKey)

	err := h.service.SendRiskNotification(ctx, in)
	if err != nil {
		// Jika service mengembalikan error, kirim Ack 'false'
		return &pb.NotificationAck{Success: false, DeliveryStatus: "FAILED"}, nil
	}

	// Jika service berhasil, kirim Ack 'true'
	return &pb.NotificationAck{Success: true, DeliveryStatus: "SENT"}, nil
}
