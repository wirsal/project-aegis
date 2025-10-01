package service

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Driver PostgreSQL
	pb "github.com/wirsal/project-aegis/api/protos"
)

// Service adalah lapisan yang menangani logika bisnis database.
type Service struct {
	db *sqlx.DB
}

// NewService membuat koneksi ke database dan mengembalikan instance Service.
func NewService(dataSourceName string) (*Service, error) {
	// Gunakan sqlx.Connect untuk membuka koneksi dan melakukan ping sekaligus
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("could not connect to DB: %w", err)
	}

	log.Println("Successfully connected to the database.")
	return &Service{db: db}, nil
}

func (s *Service) StoreRawTransaction(ctx context.Context, trx *pb.Transaction) error {
	model, err := mapToTransactionLogModel(trx)
	if err != nil {
		return fmt.Errorf("failed to map transaction to model: %w", err)
	}

	_, err = s.db.NamedExecContext(ctx, insertRawQuery, model)
	if err != nil {
		log.Printf("ERROR: failed to insert raw transaction: %v", err)
		return fmt.Errorf("failed to insert raw transaction: %w", err)
	}

	log.Printf("Successfully stored raw transaction with ID: %s", model.TrxID)
	return nil
}

func (s *Service) StoreTransaction(ctx context.Context, req *pb.StoreTransactionRequest) error {
	trx := req.GetTransactionData()
	risk := req.GetRiskData()

	// Contoh query SQL INSERT. Sesuaikan nama tabel dan kolom Anda.
	query := `INSERT INTO processed_transactions (trx_id, card_number, amount, currency, risk_level, risk_score, triggered_rules) 
	           VALUES ($1, $2, $3, $4, $5, $6, $7)`
	// println("queryStr:", query)

	_, err := s.db.ExecContext(ctx, query, trx.TrxId, trx.CardNumber, trx.TrxAmount, trx.TrxCurrency, risk.GetRiskLevel().String(), risk.GetRiskScore(), risk.GetTriggeredRules())
	if err != nil {
		return fmt.Errorf("failed to insert processed transaction: %w", err)
	}

	log.Printf("Successfully stored processed transaction with ID: %s", trx.TrxId)
	return nil
}
