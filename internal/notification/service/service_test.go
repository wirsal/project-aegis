package service

// import (
// 	"context"
// 	"errors"
// 	"io/ioutil"
// 	"net/http"
// 	"strings"
// 	"testing"

// 	pb "github.com/wirsal/project-aegis/api/protos"
// )

// // --- Mocking ---

// // mockRoundTripper adalah implementasi palsu dari http.RoundTripper.
// // Ia akan mencegat panggilan HTTP dan mengembalikan respons yang sudah kita siapkan.
// type mockRoundTripper struct {
// 	// a a a
// 	Response *http.Response
// 	Err      error
// }

// // RoundTrip adalah satu-satunya metode yang dibutuhkan oleh interface http.RoundTripper.
// func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
// 	// Kita bisa menambahkan logika di sini untuk memeriksa `req` jika perlu
// 	return m.Response, m.Err
// }

// // --- Unit Test ---

// func TestSendRiskNotification(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		mockResponse  *http.Response // Respons palsu yang akan kita kembalikan
// 		mockError     error          // Error palsu yang akan kita kembalikan
// 		inputRiskData *pb.RiskResult
// 		expectError   bool
// 	}{
// 		{
// 			name: "Success case - API returns 200 OK",
// 			mockResponse: &http.Response{
// 				StatusCode: 200,
// 				Body:       ioutil.NopCloser(strings.NewReader(`{"status":"ok"}`)),
// 				Header:     make(http.Header),
// 			},
// 			mockError: nil,
// 			inputRiskData: &pb.RiskResult{
// 				TrxKey:    "trx-123",
// 				RiskLevel: pb.RiskResult_HIGH,
// 				RiskScore: 100,
// 			},
// 			expectError: false,
// 		},
// 		{
// 			name: "Failure case - API returns 500 Internal Server Error",
// 			mockResponse: &http.Response{
// 				StatusCode: 500,
// 				Body:       ioutil.NopCloser(strings.NewReader(`{"error":"server down"}`)),
// 				Header:     make(http.Header),
// 			},
// 			mockError: nil,
// 			inputRiskData: &pb.RiskResult{
// 				TrxKey: "trx-456",
// 			},
// 			expectError: false, // Fungsi tidak mengembalikan error, hanya log
// 		},
// 		{
// 			name:         "Network error case - HTTP client fails to send",
// 			mockResponse: nil,
// 			mockError:    errors.New("network connection failed"),
// 			inputRiskData: &pb.RiskResult{
// 				TrxKey: "trx-789",
// 			},
// 			expectError: true, // Fungsi akan mengembalikan error dari http.Client
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Arrange: Buat service dengan http.Client yang sudah di-mock
// 			mockClient := &http.Client{
// 				Transport: &mockRoundTripper{
// 					Response: tt.mockResponse,
// 					Err:      tt.mockError,
// 				},
// 			}

// 			// Buat instance Service baru dan suntikkan mockClient
// 			s := NewService("https://fake-api.com/notify")
// 			s.httpClient = mockClient // Ganti httpClient asli dengan mock

// 			// Act: Panggil fungsi yang akan diuji
// 			err := s.SendRiskNotification(context.Background(), tt.inputRiskData)

// 			// Assert: Verifikasi hasilnya
// 			if (err != nil) != tt.expectError {
// 				t.Errorf("SendRiskNotification() error = %v, expectError %v", err, tt.expectError)
// 			}
// 		})
// 	}
// }
