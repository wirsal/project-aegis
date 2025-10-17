package service

import (
	"context"
	"log"
	"net/http"

	pb "github.com/wirsal/project-aegis/api/protos"
)

type WASender struct {
	httpClient *http.Client
	waUrl      string
}

func NewWASender(url string) *WASender {
	return &WASender{
		httpClient: &http.Client{},
		waUrl:      url,
	}
}

func (s *WASender) Send(ctx context.Context, req *pb.NotificationRequest) error {
	log.Printf("WASender: Preparing notification for recipient: %s", req.GetRecipient())

	return nil
}
