package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	pb "github.com/wirsal/project-aegis/api/protos"
)

// Rule mendefinisikan struktur sebuah aturan.
type Rule struct {
	RuleCode       string `json:"rule_code"`
	RuleType       string `json:"rule_type"`
	Org            string `json:"org"`
	Type           string `json:"type"`
	BlockCode      string `json:"block_code"`
	CrLimit        string `json:"cr_limit"`
	MerchCategory  string `json:"merch_category"`
	TransCode      string `json:"trans_code"`
	CountryCode    string `json:"country_code"`
	CurrencyCode   string `json:"currency_code"`
	Amount         string `json:"amount"`
	PosCondCode    string `json:"pos_cond_code"`
	RespCode       string `json:"resp_code"`
	TimeStamp      string `json:"time_stamp"`
	InstallmentInd string `json:"installment_ind"`
	FirstUsageFlag string `json:"first_usage_flag"`
	CardList       string `json:"card_list"`
}

type RuleEngineServer struct {
	pb.UnimplementedRuleEngineServer
	rules []Rule // Menyimpan semua aturan yang dimuat
}

// NewRuleEngineServer sekarang memuat aturan saat inisialisasi.
func NewRuleEngineServer(rulesPath string) (*RuleEngineServer, error) {
	rules, err := loadRules(rulesPath)
	if err != nil {
		return nil, fmt.Errorf("gagal memuat aturan: %w", err)
	}
	log.Printf("Berhasil memuat %d aturan dari %s", len(rules), rulesPath)
	return &RuleEngineServer{rules: rules}, nil
}

// loadRules membaca dan mem-parsing file JSON aturan.
func loadRules(path string) ([]Rule, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rules []Rule
	err = json.Unmarshal(file, &rules)
	return rules, err
}

func (s *RuleEngineServer) AnalyzeTransaction(ctx context.Context, in *pb.Transaction) (*pb.RiskResult, error) {
	log.Printf("Menerima transaksi untuk dianalisis: Reff=%s", in.TrxReffNumber)

	// Iterasi melalui semua aturan yang telah dimuat
	for _, rule := range s.rules {
		// Panggil fungsi validasi utama
		if s.validateRule(rule, in) {
			log.Printf("✅ Transaksi cocok dengan Aturan: %s", rule.RuleCode)

			// TODO: Panggil Persistence Service di sini jika diperlukan

			// Jika cocok, buat hasil dan langsung kembalikan
			return &pb.RiskResult{
				Rrn:            in.TrxReffNumber,
				RiskLevel:      pb.RiskResult_HIGH, // Atau tentukan dari tipe aturan
				TriggeredRules: []string{rule.RuleCode},
				RiskScore:      100, // Atau tentukan dari aturan
			}, nil
		}
	}

	log.Printf("Transaksi tidak cocok dengan aturan manapun.")
	// Jika tidak ada aturan yang cocok
	return &pb.RiskResult{
		Rrn:            in.TrxReffNumber,
		RiskLevel:      pb.RiskResult_LOW,
		TriggeredRules: []string{},
		RiskScore:      0,
	}, nil
}

// validateRule adalah implementasi Go dari fungsi JS Anda.
func (s *RuleEngineServer) validateRule(rule Rule, trx *pb.Transaction) bool {
	// Konversi amount ke int64 untuk perbandingan
	// amountInt := int64(trx.TrxAmount)
	// Konversi trx time ke int64 (hhmmss)
	// trxTimeInt, _ := strconv.ParseInt(trx.TrxTime, 10, 64)

	// Lakukan semua validasi, jika satu saja gagal, hasilnya akan false
	// return vInList(rule.Org, trx.CardOrg, "000") &&
	// 	vInList(rule.Type, trx.CardType, "000") &&
	// 	vInList(rule.MerchCategory, trx.MerchCategory, "0000") &&
	// 	vInList(rule.TransCode, trx.TrxCode, "000") &&
	// 	vCountry(rule.CountryCode, trx.TrxCountry) &&
	// 	vInRange(rule.Amount, amountInt) &&
	// 	vInRange(rule.TimeStamp, trxTimeInt) &&
	// 	vInList(rule.RespCode, trx.TrxRespCode, "AA")

	return vInList(rule.CountryCode, trx.TrxCountry, "360")
	// Tambahkan validasi lain di sini...
}
