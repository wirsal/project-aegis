package handler

import (
	"context"
	"reflect"
	"testing"

	pb "github.com/wirsal/project-aegis/api/protos"
)

// --- Mocking ---

// mockRuleEngineService adalah implementasi palsu dari interface RuleEngineService.
type mockRuleEngineService struct {
	// Tentukan hasil yang ingin dikembalikan oleh mock ini
	MockResult *pb.RiskResult
}

// Implementasikan metode AnalyzeTransaction yang dibutuhkan oleh interface.
func (m *mockRuleEngineService) AnalyzeTransaction(in *pb.Transaction) *pb.RiskResult {
	// Cukup kembalikan hasil yang sudah kita siapkan.
	return m.MockResult
}

// --- Unit Test ---

func TestAnalyzeTransaction(t *testing.T) {
	// Siapkan data input dan output yang diharapkan untuk tes
	sampleInput := &pb.Transaction{TrxKey: "trx-123"}
	expectedResult := &pb.RiskResult{
		TrxKey:    "trx-123",
		RiskLevel: pb.RiskResult_HIGH,
		RiskScore: 100,
	}

	// Buat mock service dan atur agar ia mengembalikan `expectedResult`
	mockService := &mockRuleEngineService{
		MockResult: expectedResult,
	}

	// Arrange: Buat handler dengan service palsu (mock)
	handler := NewGRPCHandler(mockService)

	// Act: Panggil metode gRPC yang akan diuji
	result, err := handler.AnalyzeTransaction(context.Background(), sampleInput)

	// Assert: Verifikasi hasilnya
	if err != nil {
		t.Fatalf("AnalyzeTransaction() returned an unexpected error: %v", err)
	}

	// Gunakan reflect.DeepEqual untuk membandingkan isi dari struct hasil
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("AnalyzeTransaction() got = %v, want %v", result, expectedResult)
	}
}
