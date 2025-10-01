package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Driver PostgreSQL
	pb "github.com/wirsal/project-aegis/api/protos"
)

// Service adalah lapisan yang menangani logika bisnis database.
type Service struct {
	db *sql.DB
}

// NewService membuat koneksi ke database dan mengembalikan instance Service.
func NewService(dataSourceName string) (*Service, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping DB: %w", err)
	}

	log.Println("Successfully connected to the database.")
	return &Service{db: db}, nil
}

// StoreRawTransaction menyimpan data transaksi mentah ke database.
func (s *Service) StoreRawTransaction(ctx context.Context, trx *pb.Transaction) error {
	// Contoh query SQL INSERT. Sesuaikan nama tabel dan kolom Anda.
	query := `INSERT INTO transaction_log (trx_id, card_acct, amount, currency, trx_date, trx_time) 
	           VALUES ($1, $2, $3, $4, $5, $6)`
	println("query:", query)
	_, err := s.db.ExecContext(ctx, query, trx.TrxId, trx.CardNumber, trx.TrxAmount, trx.TrxCurrency, trx.TrxDate, trx.TrxTime)
	if err != nil {
		return fmt.Errorf("failed to insert raw transaction: %w", err)
	}

	log.Printf("Successfully stored raw transaction with ID: %s", trx.TrxId)
	return nil
}

// StoreTransaction menyimpan data transaksi lengkap beserta hasil risikonya.
func (s *Service) StoreTransaction(ctx context.Context, req *pb.StoreTransactionRequest) error {
	trx := req.GetTransactionData()
	risk := req.GetRiskData()

	// Contoh query SQL INSERT. Sesuaikan nama tabel dan kolom Anda.
	query := `INSERT INTO processed_transactions (trx_id, card_number, amount, currency, risk_level, risk_score, triggered_rules) 
	           VALUES ($1, $2, $3, $4, $5, $6, $7)`
	println("queryStr:", query)

	_, err := s.db.ExecContext(ctx, query, trx.TrxId, trx.CardNumber, trx.TrxAmount, trx.TrxCurrency, risk.GetRiskLevel().String(), risk.GetRiskScore(), risk.GetTriggeredRules())
	if err != nil {
		return fmt.Errorf("failed to insert processed transaction: %w", err)
	}

	log.Printf("Successfully stored processed transaction with ID: %s", trx.TrxId)
	return nil
}
