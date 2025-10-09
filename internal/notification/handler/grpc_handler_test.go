package handler

import (
	"context"
	"errors"
	"testing"

	pb "github.com/wirsal/project-aegis/api/protos"
)

// --- Mocking ---

// mockNotificationService adalah implementasi palsu dari interface NotificationService.
type mockNotificationService struct {
	// Kita bisa menambahkan field untuk mengontrol perilakunya
	shouldReturnError bool
	wasCalled         bool
}

// Implementasikan metode yang dibutuhkan oleh interface.
func (m *mockNotificationService) SendRiskNotification(ctx context.Context, riskData *pb.RiskResult) error {
	m.wasCalled = true // Tandai bahwa metode ini telah dipanggil
	if m.shouldReturnError {
		return errors.New("simulated service error")
	}
	return nil
}

// --- Unit Test ---

func TestSendRiskNotification(t *testing.T) {
	// Definisikan kasus-kasus uji
	tests := []struct {
		name              string
		mockService       *mockNotificationService
		input             *pb.RiskResult
		expectedSuccess   bool
		expectedStatus    string
		expectServiceCall bool
	}{
		{
			name: "Success case - service returns no error",
			mockService: &mockNotificationService{
				shouldReturnError: false,
			},
			input:             &pb.RiskResult{TrxKey: "rrn-123"},
			expectedSuccess:   true,
			expectedStatus:    "SENT",
			expectServiceCall: true,
		},
		{
			name: "Failure case - service returns an error",
			mockService: &mockNotificationService{
				shouldReturnError: true,
			},
			input:             &pb.RiskResult{TrxKey: "rrn-456"},
			expectedSuccess:   false,
			expectedStatus:    "FAILED",
			expectServiceCall: true,
		},
	}

	// Jalankan setiap kasus uji
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: Buat handler dengan service palsu (mock)
			handler := NewGRPCHandler(tt.mockService)

			// Act: Panggil metode gRPC yang akan diuji
			ack, err := handler.SendRiskNotification(context.Background(), tt.input)

			// Assert: Verifikasi hasilnya
			if err != nil {
				t.Errorf("SendRiskNotification() returned an unexpected error: %v", err)
			}

			if ack.Success != tt.expectedSuccess {
				t.Errorf("Expected Ack.Success to be %v, but got %v", tt.expectedSuccess, ack.Success)
			}

			if ack.DeliveryStatus != tt.expectedStatus {
				t.Errorf("Expected Ack.DeliveryStatus to be '%s', but got '%s'", tt.expectedStatus, ack.DeliveryStatus)
			}

			if tt.mockService.wasCalled != tt.expectServiceCall {
				t.Errorf("Expected service call to be %v, but was %v", tt.expectServiceCall, tt.mockService.wasCalled)
			}
		})
	}
}
