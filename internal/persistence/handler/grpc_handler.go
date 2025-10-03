package handler

import (
	"context"
	"log"

	pb "github.com/wirsal/project-aegis/api/protos"
	"github.com/wirsal/project-aegis/internal/persistence/service"
)

// GRPCHandler mengimplementasikan interface gRPC Server.
type GRPCHandler struct {
	pb.UnimplementedPersistenceServer
	service *service.Service
}

// NewGRPCHandler membuat instance handler baru.
func NewGRPCHandler(svc *service.Service) *GRPCHandler {
	return &GRPCHandler{
		service: svc,
	}
}

// StoreRawTransaction hanya mendelegasikan tugas ke service layer.
func (h *GRPCHandler) StoreRawTransaction(ctx context.Context, in *pb.Transaction) (*pb.StoreAck, error) {
	log.Printf("gRPC handler received StoreRawTransaction for TrxKey: %s", in.TrxKey)
	err := h.service.StoreRawTransaction(ctx, in)
	if err != nil {
		return &pb.StoreAck{Success: false, Message: err.Error()}, nil
	}
	return &pb.StoreAck{Success: true, Message: "Raw transaction stored successfully"}, nil
}

// StoreTransaction hanya mendelegasikan tugas ke service layer.
func (h *GRPCHandler) StoreTransaction(ctx context.Context, in *pb.StoreTransactionRequest) (*pb.StoreAck, error) {
	log.Printf("gRPC handler received StoreTransaction for TrxKey: %s", in.GetTransactionData().GetTrxKey())
	err := h.service.StoreRiskResult(ctx, in)
	if err != nil {
		return &pb.StoreAck{Success: false, Message: err.Error()}, nil
	}
	return &pb.StoreAck{Success: true, Message: "Processed transaction stored successfully"}, nil
}
