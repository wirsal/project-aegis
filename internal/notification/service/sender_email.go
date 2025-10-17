package service

import (
	"context"
	"log"
	"net/http"

	pb "github.com/wirsal/project-aegis/api/protos"
)

type EmailSender struct {
	httpClient *http.Client
	emailUrl   string
}

func NewEmailSender(url string) *EmailSender {
	return &EmailSender{
		httpClient: &http.Client{},
		emailUrl:   url,
	}
}

func (s *EmailSender) Send(ctx context.Context, req *pb.NotificationRequest) error {
	log.Printf("EmailSender: Preparing notification for recipient: %s", req.GetRecipient())

	return nil
}
