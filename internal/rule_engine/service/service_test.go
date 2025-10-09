package service

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"

	pb "github.com/wirsal/project-aegis/api/protos"
)

// --- Mocking Dependencies ---
type mockPersistenceClient struct {
	storeTransactionCalled bool
	receivedTrx            *pb.Transaction
	receivedRisk           *pb.RiskResult
}

func (m *mockPersistenceClient) StoreTransaction(ctx context.Context, in *pb.StoreTransactionRequest, opts ...grpc.CallOption) (*pb.StoreAck, error) {
	m.storeTransactionCalled = true
	m.receivedTrx = in.GetTransactionData()
	m.receivedRisk = in.GetRiskData()
	return &pb.StoreAck{Success: true}, nil
}

type mockNotificationClient struct {
	called       bool
	receivedRisk *pb.RiskResult
}

func (m *mockNotificationClient) SendRiskNotification(ctx context.Context, in *pb.RiskResult, opts ...grpc.CallOption) (*pb.NotificationAck, error) {
	m.called = true
	m.receivedRisk = in
	return &pb.NotificationAck{Success: true}, nil
}

func (m *mockPersistenceClient) StoreRawTransaction(ctx context.Context, in *pb.Transaction, opts ...grpc.CallOption) (*pb.StoreAck, error) {
	return &pb.StoreAck{Success: true}, nil
}

func TestAnalyzeTransaction(t *testing.T) {
	defaultValidRule := Rule{
		Org:            "000",
		Type:           "000",
		CountryCode:    "A000",
		CurrencyCode:   "A000",
		MerchCategory:  "0000",
		PosCondCode:    "AA",
		RespCode:       "AA",
		TimeStamp:      "000000-235959",
		InstallmentInd: "-",
	}

	tests := []struct {
		name                   string
		rules                  []Rule
		inputTrx               *pb.Transaction
		expectedRiskLevel      pb.RiskResult_RiskLevel
		expectMatch            bool
		expectedRuleCode       string
		expectPersistenceCall  bool
		expectNotificationCall bool
	}{
		{
			name: "Transaction that matches a rule on Amount",
			rules: func() []Rule {
				rule := defaultValidRule
				rule.Priority = 1
				rule.RuleCode = "HIGH_AMOUNT"
				rule.Amount = "1000-5000"
				return []Rule{rule}
			}(),
			inputTrx: &pb.Transaction{
				TrxKey:         "key-123",
				TrxAmount:      2500,
				CardOrg:        "001",
				CardType:       "002",
				TrxCountry:     "360",
				TrxCurrency:    "360",
				MerchCategory:  "1234",
				TrxPosMode:     "01",
				TrxRespCode:    "00",
				TrxTime:        "123000",
				TrxInstallment: "N",
			},
			expectedRiskLevel:      pb.RiskResult_HIGH,
			expectMatch:            true,
			expectedRuleCode:       "HIGH_AMOUNT",
			expectPersistenceCall:  true,
			expectNotificationCall: true,
		},
		{
			name: "Transaction that does not match any rule",
			rules: func() []Rule {
				rule := defaultValidRule
				rule.Priority = 1
				rule.RuleCode = "HIGH_AMOUNT"
				rule.Amount = "1000-5000"
				return []Rule{rule}
			}(),
			inputTrx: &pb.Transaction{
				TrxKey:         "key-456",
				TrxAmount:      99,
				CardOrg:        "001",
				CardType:       "002",
				TrxCountry:     "360",
				TrxCurrency:    "360",
				MerchCategory:  "1234",
				TrxPosMode:     "01",
				TrxRespCode:    "00",
				TrxTime:        "123000",
				TrxInstallment: "N",
			},
			expectedRiskLevel:      pb.RiskResult_LOW,
			expectMatch:            false,
			expectPersistenceCall:  false,
			expectNotificationCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPsc := &mockPersistenceClient{}
			mockNsc := &mockNotificationClient{}

			s := &Service{
				rules:              tt.rules,
				persistenceClient:  mockPsc,
				notificationClient: mockNsc,
			}

			result := s.AnalyzeTransaction(tt.inputTrx)
			time.Sleep(10 * time.Millisecond)

			if result.RiskLevel != tt.expectedRiskLevel {
				t.Errorf("Expected risk level %v, but got %v", tt.expectedRiskLevel, result.RiskLevel)
			}

			if tt.expectMatch {
				if len(result.TriggeredRules) == 0 {
					t.Fatal("Expected a rule to be triggered, but none was.")
				}
				if result.TriggeredRules[0] != tt.expectedRuleCode {
					t.Errorf("Expected triggered rule '%s', but got '%s'", tt.expectedRuleCode, result.TriggeredRules[0])
				}
			} else {
				if len(result.TriggeredRules) > 0 {
					t.Errorf("Expected no rules to be triggered, but got %v", result.TriggeredRules)
				}
			}

			if mockPsc.storeTransactionCalled != tt.expectPersistenceCall {
				t.Errorf("Expected persistence client call to be %v, but was %v", tt.expectPersistenceCall, mockPsc.storeTransactionCalled)
			}
			if mockNsc.called != tt.expectNotificationCall {
				t.Errorf("Expected notification client call to be %v, but was %v", tt.expectNotificationCall, mockNsc.called)
			}
		})
	}
}
