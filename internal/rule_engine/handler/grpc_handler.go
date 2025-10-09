package handler

import (
	"context"

	pb "github.com/wirsal/project-aegis/api/protos"
)

type RuleEngineService interface {
	AnalyzeTransaction(in *pb.Transaction) *pb.RiskResult
}

type GRPCHandler struct {
	pb.UnimplementedRuleEngineServer
	service RuleEngineService
}

func NewGRPCHandler(svc RuleEngineService) *GRPCHandler {
	return &GRPCHandler{
		service: svc,
	}
}

func (h *GRPCHandler) AnalyzeTransaction(ctx context.Context, in *pb.Transaction) (*pb.RiskResult, error) {
	result := h.service.AnalyzeTransaction(in)
	return result, nil
}
