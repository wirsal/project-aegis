package service

import (
	"context"
	"fmt"
	"log"

	pb "github.com/wirsal/project-aegis/api/protos"
	"github.com/wirsal/project-aegis/pkg/config"
	"google.golang.org/protobuf/types/known/structpb"
)

type Sender interface {
	Send(ctx context.Context, req *pb.NotificationRequest) error
}

type Service struct {
	senders map[string]Sender
}

func NewService(cfg config.NotificationConfig) *Service {
	firebaseSender := NewFirebaseSender(cfg.FCMgatewayURL)
	// smsSender := NewSMSSender(cfg.SMSGatewayURL)

	return &Service{
		senders: map[string]Sender{
			"FIREBASE": firebaseSender,
			// "SMS":      smsSender,
		},
	}
}

func (s *Service) TriggerRiskAlert(ctx context.Context, req *pb.RiskAlertRequest) error {
	trxData := req.GetTransactionData()
	riskData := req.GetRiskData()

	// 1. Tentukan Channel (contoh logika sederhana)
	//    Logika ini bisa lebih kompleks, misal membaca dari config atau database.
	var channel string
	if riskData.GetRiskScore() > 90 {
		channel = "FIREBASE" // atau "SMS"
	} else {
		channel = "FIREBASE"
	}

	// 2. Cari Penerima (ini bagian penting yang perlu Anda kembangkan)
	//    TODO: Implementasikan logic untuk mencari device token dari database pelanggan
	//          berdasarkan trxData.GetCardNumber() atau ID pelanggan.
	recipient := "ch30ZNhq7kDgnPqChvHg6W:APA91bGn9oH5JuUh4BBV2gF_0B4dZXLfF2Yo94xFX7dr_w5awHjEZBfTpOjg1IwI1F-C6Ap9KB6lZGb2tFet4W3jEvviJ6q1aMzI9pNu_4iaqbpj12-FNfw" // Placeholder
	if recipient == "" {
		return fmt.Errorf("recipient not found for card number %s", trxData.GetCardNumber())
	}

	// 3. Susun Payload Pesan
	payloadData := map[string]interface{}{
		"title": "Peringatan Keamanan Kartu Anda",
		"body":  fmt.Sprintf("Terdeteksi transaksi mencurigakan sebesar %.2f di %s.", trxData.GetTrxAmount(), trxData.GetTrxMerchantName()),
	}
	payloadStruct, err := structpb.NewStruct(payloadData)
	if err != nil {
		return fmt.Errorf("failed to create payload struct: %w", err)
	}

	// 4. Buat Request Internal untuk didistribusikan ke Sender
	internalReq := &pb.NotificationRequest{
		Channel:   channel,
		Recipient: recipient,
		Payload:   payloadStruct,
	}

	// 5. Panggil distributor internal (yang sudah ada sebelumnya)
	return s.SendNotification(ctx, internalReq)
}

func (s *Service) SendNotification(ctx context.Context, req *pb.NotificationRequest) error {
	channel := req.GetChannel()
	log.Printf("Distributing notification for channel: %s", channel)

	sender, ok := s.senders[channel]
	if !ok {
		return fmt.Errorf("unsupported notification channel: %s", channel)
	}

	return sender.Send(ctx, req)
}
