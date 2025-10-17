package service

import (
	"context"
	"fmt"
	"log"
	"strings"

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

type Recipient struct {
	fcm_token    string
	email        string
	phone_number string
}

func NewService(cfg config.NotificationConfig) *Service {
	return &Service{
		senders: map[string]Sender{
			"firebase": NewFirebaseSender(cfg.FCMgatewayURL),
			"email":    NewEmailSender(cfg.EmailGatewayURL),
			"sms":      NewWASender(cfg.WAGatewayURL),
			"wa":       NewWASender(cfg.WAGatewayURL),
		},
	}
}

func (s *Service) TriggerRiskAlert(ctx context.Context, req *pb.RiskAlertRequest) error {
	trxData := req.GetTransactionData()
	riskData := req.GetRiskData()

	recipientInfo := getRecipient(trxData.CardNumber)

	channels := strings.Split(riskData.RuleChannel, ",")
	var processingErrors []string // Slice untuk mengumpulkan pesan error

	for _, channel := range channels {
		cleanedChannel := strings.TrimSpace(strings.ToLower(channel))
		if cleanedChannel == "" {
			continue
		}

		var recipientAddress string
		var payloadStruct *structpb.Struct
		switch cleanedChannel {
		case "firebase":
			recipientAddress = recipientInfo.fcm_token
			payloadStruct, _ = getPayloadFCM(trxData)
		case "email":
			recipientAddress = recipientInfo.email
			payloadStruct, _ = getPayloadEmail(trxData, riskData.RuleTemplatesId)
		case "wa", "sms":
			recipientAddress = recipientInfo.phone_number
			payloadStruct, _ = getPayloadWa(trxData, riskData.RuleTemplatesId)
		default:
			log.Printf("Warning: channel '%s' is not supported, skipping.", cleanedChannel)
			continue
		}

		if recipientAddress == "" {
			log.Printf("Warning: recipient address for channel '%s' is empty, skipping.", cleanedChannel)
			continue
		}

		internalReq := &pb.NotificationRequest{
			Channel:   cleanedChannel,
			Recipient: recipientAddress,
			Payload:   payloadStruct,
		}

		log.Printf("Attempting to send notification via channel: %s", cleanedChannel)
		if err := s.SendNotification(ctx, internalReq); err != nil {
			processingErrors = append(processingErrors, fmt.Sprintf("channel %s: %v", cleanedChannel, err))
		}
	}

	if len(processingErrors) > 0 {
		return fmt.Errorf("encountered errors during notification processing: %s", strings.Join(processingErrors, "; "))
	}

	log.Printf("Successfully processed all notification channels for TrxKey: %s", trxData.GetTrxKey())
	return nil
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
