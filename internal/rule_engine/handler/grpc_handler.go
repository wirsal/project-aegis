package handler

import (
	"context"

	pb "github.com/wirsal/project-aegis/api/protos"
	"github.com/wirsal/project-aegis/internal/rule_engine/service"
)

// GRPCHandler adalah lapisan yang menerjemahkan panggilan gRPC ke logika bisnis.
type GRPCHandler struct {
	pb.UnimplementedRuleEngineServer
	service *service.Service // Dependensi ke service layer
}

// NewGRPCHandler membuat instance handler baru.
func NewGRPCHandler(svc *service.Service) *GRPCHandler {
	return &GRPCHandler{
		service: svc,
	}
}

// AnalyzeTransaction sekarang menjadi method yang sangat tipis.
// Tugasnya hanya menerima panggilan gRPC dan mendelegasikannya ke service.
func (h *GRPCHandler) AnalyzeTransaction(ctx context.Context, in *pb.Transaction) (*pb.RiskResult, error) {
	// Panggil logika bisnis yang sebenarnya dari service layer
	result := h.service.AnalyzeTransaction(in)

	// Kembalikan hasilnya. Error di sini nil karena service kita tidak mengembalikan error.
	return result, nil
}
