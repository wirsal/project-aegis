package service

import (
	"context"
	"log"

	"google.golang.org/grpc"

	pb "github.com/wirsal/project-aegis/api/protos"
	"github.com/wirsal/project-aegis/pkg/parser"
)

type RuleEngineClient interface {
	AnalyzeTransaction(ctx context.Context, in *pb.Transaction, opts ...grpc.CallOption) (*pb.RiskResult, error)
}

type PersistenceClient interface {
	StoreRawTransaction(ctx context.Context, in *pb.Transaction, opts ...grpc.CallOption) (*pb.StoreAck, error)
}

type GatewayService struct {
	ruleEngineClient  RuleEngineClient
	persistenceClient PersistenceClient
}

func NewGatewayService(reClient RuleEngineClient, psClient PersistenceClient) *GatewayService {
	return &GatewayService{
		ruleEngineClient:  reClient,
		persistenceClient: psClient,
	}
}

// ProcessAndForwardMessage sekarang lebih bersih dan fokus pada alur kerja
func (s *GatewayService) ProcessAndForwardMessage(ctx context.Context, rawMessage []byte) error {
	log.Println("Passing raw message to parser...")

	// 1. Panggil parser dari paket terpisah
	protoMsg, err := parser.ParseAndMapTransaction(string(rawMessage))
	if err != nil {
		log.Printf("ERROR: Failed to parse transaction: %v", err)
		return err
	}

	// 2. Panggil Persistence Service untuk menyimpan transaksi mentah (secara asinkron)
	go s.storeRawTransaction(protoMsg)

	// 3. Panggil Rule Engine Service untuk analisis
	log.Printf("✅ Message parsed. Sending transaction (ID: %s) to Rule Engine...", protoMsg.TrxKey)
	riskResult, err := s.ruleEngineClient.AnalyzeTransaction(ctx, protoMsg)
	if err != nil {
		log.Printf("ERROR: gRPC call to Rule Engine failed: %v", err)
		return err
	}

	log.Printf("✅ Response from Rule Engine received for TrxKey %s. Risk Score: %d", protoMsg.TrxKey, riskResult.RiskScore)
	return nil
}

func (s *GatewayService) storeRawTransaction(trx *pb.Transaction) {
	log.Printf("Storing raw transaction (TrxKey: %s) to persistence...", trx.TrxKey)
	_, err := s.persistenceClient.StoreRawTransaction(context.Background(), trx)
	if err != nil {
		log.Printf("ERROR: Failed to store raw transaction for TrxKey %s: %v", trx.TrxKey, err)
	} else {
		log.Printf("✅ Raw transaction (TrxKey: %s) stored successfully.", trx.TrxKey)
	}
}
