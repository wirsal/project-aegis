package service

import (
	"reflect" // Diperlukan untuk membandingkan struct secara mendalam
	"testing"
	"time"

	pb "github.com/wirsal/project-aegis/api/protos"
)

func TestSafeAtoi(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"Valid Positive", "123", 123},
		{"Valid Negative", "-45", -45},
		{"Zero", "0", 0},
		{"Invalid String", "abc", 0},
		{"Empty String", "", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := safeAtoi(tt.input); got != tt.want {
				t.Errorf("safeAtoi() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeParseFloat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  float64
	}{
		{"Valid Float", "123.45", 123.45},
		{"Valid Integer", "789", 789.0},
		{"Invalid String", "xyz", 0.0},
		{"Empty String", "", 0.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := safeParseFloat(tt.input); got != tt.want {
				t.Errorf("safeParseFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapToTransactionLogModel(t *testing.T) {
	sampleTrx := &pb.Transaction{
		TrxKey:        "test-key-123",
		CardOrg:       "001",
		CardType:      "002",
		CardNumber:    "4000123456789010",
		CardExpired:   "1228",
		TrxDate:       "2025-10-09",
		TrxTime:       "10:30:00",
		MerchOrg:      "003",
		MerchNumber:   "merchant001",
		TrxCardType:   "C",
		TrxCode:       "01",
		TrxRespCode:   "00",
		TrxAmount:     15000.50,
		TrxBillAmount: "15000.50",
		MerchCategory: "5411", // Grocery
		TrxCountry:    "360",
		TrxChipData:   "somechipdata",
	}

	expectedTime, _ := time.Parse("2006-01-02 15:04:05", "2025-10-09 10:30:00")
	expectedModel := &TransactionLogModel{
		TrxKey:         "test-key-123",
		CardOrg:        "001",
		CardType:       "002",
		CardNumber:     "4000123456789010",
		CardExpDate:    "1228",
		TrxDate:        "2025-10-09",
		TrxTime:        "10:30:00",
		TrxDatetime:    expectedTime,
		MerchOrg:       "003",
		MerchID:        "merchant001",
		TrxCardType:    "C",
		TrxCode:        1,
		TrxRespCode:    "00",
		TrxBillAmt:     15000.50,
		TrxAmt:         15000.50,
		TrxMCC:         5411,
		TrxCountryCode: 360,
		TrxChipLength:  12,
		TrxChipData:    "somechipdata",
	}

	t.Run("Happy Path - All data valid", func(t *testing.T) {
		got, err := mapToTransactionLogModel(sampleTrx)
		if err != nil {
			t.Fatalf("mapToTransactionLogModel() returned an unexpected error: %v", err)
		}

		if !reflect.DeepEqual(got, expectedModel) {
			t.Errorf("mapToTransactionLogModel() got = \n%v, \nwant \n%v", got, expectedModel)
		}
	})

	t.Run("Invalid numbers should be converted to zero", func(t *testing.T) {
		invalidTrx := &pb.Transaction{
			TrxCode:       "abc",
			MerchCategory: "xyz",
		}

		got, _ := mapToTransactionLogModel(invalidTrx)

		if got.TrxCode != 0 {
			t.Errorf("Expected TrxCode to be 0 for invalid input, but got %d", got.TrxCode)
		}
		if got.TrxMCC != 0 {
			t.Errorf("Expected TrxMCC to be 0 for invalid input, but got %d", got.TrxMCC)
		}
	})
}
