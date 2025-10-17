package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"

	pb "github.com/wirsal/project-aegis/api/protos"
)

type Service struct {
	rules              []Rule
	persistenceClient  pb.PersistenceClient
	notificationClient pb.NotificationClient
}

// Ganti nama constructor menjadi 'NewService'
func NewService(rulesPath string, psc pb.PersistenceClient, nsc pb.NotificationClient) (*Service, error) {
	rules, err := loadRules(rulesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load rules: %w", err)
	}

	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority < rules[j].Priority
	})

	log.Printf("Successfully loaded and sorted %d rules from %s", len(rules), rulesPath)
	return &Service{
		rules:              rules,
		persistenceClient:  psc,
		notificationClient: nsc,
	}, nil
}

func loadRules(path string) ([]Rule, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rules []Rule
	err = json.Unmarshal(file, &rules)

	return rules, err
}

func (s *Service) AnalyzeTransaction(in *pb.Transaction) *pb.RiskResult {
	log.Printf("🔎 Analyzing transaction (Reff=%s) against %d rules...", in.TrxKey, len(s.rules))

	for _, rule := range s.rules {
		if s.validateRule(rule, in) {
			log.Printf("✅ Transaction MATCHED rule with priority %d: %s", rule.Priority, rule.RuleCode)

			// Buat RiskResult yang akan dikirim
			riskResult := &pb.RiskResult{
				TrxKey:         in.TrxKey,
				RiskLevel:      pb.RiskResult_HIGH,
				TriggeredRules: []string{rule.RuleCode},
				RiskScore:      100,
				RuleCode:       rule.RuleCode,
				RuleType:       rule.RuleType,
			}

			go s.callPersistence(in, riskResult)
			go s.callNotification(in, riskResult)

			return riskResult
		}
	}

	log.Printf("Transaction did not match any rules.")
	return &pb.RiskResult{
		TrxKey:    in.TrxKey,
		RiskLevel: pb.RiskResult_LOW,
	}
}

func (s *Service) callPersistence(trxData *pb.Transaction, riskData *pb.RiskResult) {
	log.Printf("Calling Persistence Service for TrxKey: %s", trxData.TrxKey)

	req := &pb.StoreTransactionRequest{
		TransactionData: trxData,
		RiskData:        riskData,
	}

	_, err := s.persistenceClient.StoreTransaction(context.Background(), req)
	if err != nil {
		log.Printf("ERROR: Failed to call Persistence Service: %v", err)
	} else {
		log.Printf("Successfully called Persistence Service for TrxKey: %s", trxData.TrxKey)
	}
}

func (s *Service) callNotification(trxData *pb.Transaction, riskData *pb.RiskResult) {
	log.Printf("Calling Notification Service to trigger alert for TrxKey: %s", riskData.TrxKey)

	alertReq := &pb.RiskAlertRequest{
		TransactionData: trxData,
		RiskData:        riskData,
	}

	_, err := s.notificationClient.TriggerRiskAlert(context.Background(), alertReq)
	if err != nil {
		log.Printf("ERROR: Failed to call Notification Service: %v", err)
	} else {
		log.Printf("Successfully triggered notification for TrxKey: %s", riskData.TrxKey)
	}
}
