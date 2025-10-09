package handler

import (
	"context"
	"errors"
	"testing"

	pb "github.com/wirsal/project-aegis/api/protos"
)

// --- Mocking ---

// mockPersistenceService adalah implementasi palsu dari interface PersistenceService.
type mockPersistenceService struct {
	shouldReturnError bool
}

func (m *mockPersistenceService) StoreRawTransaction(ctx context.Context, trx *pb.Transaction) error {
	if m.shouldReturnError {
		return errors.New("failed to store raw tx")
	}
	return nil
}

func (m *mockPersistenceService) StoreRiskResult(ctx context.Context, req *pb.StoreTransactionRequest) error {
	if m.shouldReturnError {
		return errors.New("failed to store risk result")
	}
	return nil
}

// --- Unit Tests ---

func TestStoreRawTransaction(t *testing.T) {
	tests := []struct {
		name            string
		mockService     *mockPersistenceService
		input           *pb.Transaction
		expectedSuccess bool
	}{
		{
			name:            "Success case",
			mockService:     &mockPersistenceService{shouldReturnError: false},
			input:           &pb.Transaction{TrxKey: "raw-123"},
			expectedSuccess: true,
		},
		{
			name:            "Failure case",
			mockService:     &mockPersistenceService{shouldReturnError: true},
			input:           &pb.Transaction{TrxKey: "raw-456"},
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewGRPCHandler(tt.mockService)
			ack, err := handler.StoreRawTransaction(context.Background(), tt.input)

			if err != nil {
				t.Fatalf("StoreRawTransaction() returned an unexpected error: %v", err)
			}
			if ack.Success != tt.expectedSuccess {
				t.Errorf("Expected Ack.Success to be %v, but got %v", tt.expectedSuccess, ack.Success)
			}
		})
	}
}

func TestStoreTransaction(t *testing.T) {
	tests := []struct {
		name            string
		mockService     *mockPersistenceService
		input           *pb.StoreTransactionRequest
		expectedSuccess bool
	}{
		{
			name:        "Success case",
			mockService: &mockPersistenceService{shouldReturnError: false},
			input: &pb.StoreTransactionRequest{
				TransactionData: &pb.Transaction{TrxKey: "processed-123"},
			},
			expectedSuccess: true,
		},
		{
			name:        "Failure case",
			mockService: &mockPersistenceService{shouldReturnError: true},
			input: &pb.StoreTransactionRequest{
				TransactionData: &pb.Transaction{TrxKey: "processed-456"},
			},
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewGRPCHandler(tt.mockService)
			ack, err := handler.StoreTransaction(context.Background(), tt.input)

			if err != nil {
				t.Fatalf("StoreTransaction() returned an unexpected error: %v", err)
			}
			if ack.Success != tt.expectedSuccess {
				t.Errorf("Expected Ack.Success to be %v, but got %v", tt.expectedSuccess, ack.Success)
			}
		})
	}
}
