package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"

	pb "github.com/wirsal/project-aegis/api/protos"
)

// Struct Rule tetap sama
type Rule struct {
	Priority       int    `json:"priority"`
	Status         int    `json:"status"`
	RuleCode       string `json:"rule_code"`
	RuleDesc       string `json:"rule_desc"`
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

// Ganti nama struct menjadi 'Service' untuk merepresentasikan lapisan logika bisnis
type Service struct {
	rules []Rule
}

// Ganti nama constructor menjadi 'NewService'
func NewService(rulesPath string) (*Service, error) {
	rules, err := loadRules(rulesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load rules: %w", err)
	}

	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority < rules[j].Priority
	})

	log.Printf("Successfully loaded and sorted %d rules from %s", len(rules), rulesPath)
	return &Service{rules: rules}, nil
}

// fungsi loadRules tetap sama
func loadRules(path string) ([]Rule, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rules []Rule
	err = json.Unmarshal(file, &rules)
	return rules, err
}

// AnalyzeTransaction sekarang menjadi fungsi bisnis murni.
// Perhatikan bahwa 'context' sudah tidak ada, karena itu urusan handler.
func (s *Service) AnalyzeTransaction(in *pb.Transaction) *pb.RiskResult {
	log.Printf("Analyzing transaction (Reff=%s) against %d rules...", in.TrxReffNumber, len(s.rules))

	for _, rule := range s.rules {
		if s.validateRule(rule, in) {
			log.Printf("✅ Transaction MATCHED rule with priority %d: %s %s", rule.Priority, rule.RuleCode, rule.RuleDesc)

			// TODO: Panggil Persistence Service di sini jika diperlukan

			return &pb.RiskResult{
				Rrn:            in.TrxReffNumber,
				RiskLevel:      pb.RiskResult_HIGH,
				TriggeredRules: []string{rule.RuleCode},
				RiskScore:      100,
			}
		}
	}

	log.Printf("Transaction did not match any rules.")
	return &pb.RiskResult{
		Rrn:            in.TrxReffNumber,
		RiskLevel:      pb.RiskResult_LOW,
		TriggeredRules: []string{},
		RiskScore:      0,
	}
}

// validateRule tetap sama
func (s *Service) validateRule(rule Rule, trx *pb.Transaction) bool {
	amountInt := int64(trx.TrxAmount)

	return vInList(rule.Org, trx.CardOrg, "000") &&
		vInclusionExclusion(rule.CountryCode, trx.TrxCountry, "A000") &&
		vInList(rule.Type, trx.CardType, "000") &&
		vInList(rule.MerchCategory, trx.MerchCategory, "0000") &&
		vInList(rule.PosCondCode, trx.TrxPosMode, "00") &&
		vInList(rule.RespCode, trx.TrxRespCode, "00") &&
		vInRange(rule.Amount, amountInt)
}
