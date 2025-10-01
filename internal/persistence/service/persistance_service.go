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

	log.Printf("Successfully stored raw transaction with ID: %s", model.TrxKey)
	return nil
}

func (s *Service) StoreRiskResult(ctx context.Context, req *pb.StoreTransactionRequest) error {
	model, err := mapToRiskResultModel(req)
	if err != nil {
		return fmt.Errorf("failed to map risk result to model: %w", err)
	}

	_, err = s.db.NamedExecContext(ctx, insertRiskResultQuery, model)
	if err != nil {
		log.Printf("ERROR: failed to insert risk result : %v", err)
		return fmt.Errorf("failed to insert risk result: %w", err)
	}

	log.Printf("Successfully stored risk result with ID: %s", req.TransactionData.TrxKey)
	return nil
}
