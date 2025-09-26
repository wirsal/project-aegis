package service

import (
	"context"
	"log"
	"strings"

	pb "github.com/wirsal/project-aegis/api/protos"
)

// RuleEngineServer adalah implementasi dari gRPC server kita.
// Ia harus memenuhi interface `pb.RuleEngineServer` yang dibuat oleh protoc.
type RuleEngineServer struct {
	pb.UnimplementedRuleEngineServer // Wajib untuk forward compatibility
}

// NewRuleEngineServer membuat instance baru dari server.
func NewRuleEngineServer() *RuleEngineServer {
	return &RuleEngineServer{}
}

// AnalyzeTransaction adalah implementasi dari metode RPC.
// Di sinilah semua aturan risiko (business rules) akan dievaluasi.
func (s *RuleEngineServer) AnalyzeTransaction(ctx context.Context, in *pb.Transaction) (*pb.RiskResult, error) {
	log.Printf("Menerima transaksi untuk dianalisis: Reff=%s, Amount=%.2f, CardNum=%s", in.TrxReffNumber, in.TrxAmount, in.CardNumber)

	var triggeredRules []string
	var riskScore int32 = 0

	// --- CONTOH ATURAN-ATURAN RISIKO ---

	// Aturan 1: Jumlah transaksi sangat tinggi
	if in.TrxAmount > 10000000 { // Di atas 10 juta
		riskScore += 50
		triggeredRules = append(triggeredRules, "High Transaction Amount (> 10,000,000)")
	}

	// Aturan 2: Transaksi di negara berisiko tinggi (contoh)
	// Anda bisa memiliki daftar negara berisiko di database atau file konfigurasi.
	riskyCountries := map[string]bool{"840": true, "528": true} // Contoh: USA, Netherlands
	if _, isRisky := riskyCountries[in.TrxCountry]; isRisky {
		riskScore += 30
		triggeredRules = append(triggeredRules, "Transaction from High-Risk Country")
	}

	// Aturan 3: Kartu tidak menggunakan chip (fallback atau manual entry)
	// '01' = Manual Key-in, '07' = Fallback magnetic stripe
	if in.TrxPosMode == "01" || in.TrxPosMode == "07" {
		riskScore += 25
		triggeredRules = append(triggeredRules, "Non-Chip Transaction (Manual/Fallback)")
	}

	// Aturan 4: Nama merchant mengandung kata-kata mencurigakan
	suspiciousKeywords := []string{"GAMBLING", "CASINO", "BETTING"}
	for _, keyword := range suspiciousKeywords {
		if strings.Contains(strings.ToUpper(in.TrxMerchantName), keyword) {
			riskScore += 40
			triggeredRules = append(triggeredRules, "Suspicious Merchant Keyword: "+keyword)
			break // Cukup temukan satu
		}
	}

	// Aturan 5: Transaksi ditolak oleh penerbit kartu (issuer)
	// Kode '05' atau '51' adalah contoh umum untuk Decline / Insufficient Funds
	if in.TrxRespCode == "05" || in.TrxRespCode == "51" {
		riskScore += 15
		triggeredRules = append(triggeredRules, "Transaction Declined by Issuer")
	}

	// Tentukan level risiko berdasarkan skor total
	var riskLevel pb.RiskResult_RiskLevel
	if riskScore >= 80 {
		riskLevel = pb.RiskResult_HIGH
	} else if riskScore >= 40 {
		riskLevel = pb.RiskResult_MEDIUM
	} else {
		riskLevel = pb.RiskResult_LOW
	}

	log.Printf("Analisis selesai. Skor: %d, Level: %s, Aturan Terpicu: %v", riskScore, riskLevel, triggeredRules)

	// Buat dan kirim kembali hasil analisis risiko
	return &pb.RiskResult{
		Rrn:            in.TrxReffNumber,
		RiskLevel:      riskLevel,
		TriggeredRules: triggeredRules,
		RiskScore:      riskScore,
	}, nil
}
