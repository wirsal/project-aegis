package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	pb "github.com/wirsal/project-aegis/api/protos"
)

// Service handles the business logic for sending notifications.
type Service struct {
	httpClient  *http.Client
	externalURL string // URL of the external RESTful API
}

// NewService creates a new notification service instance.
func NewService(externalAPIURL string) *Service {
	return &Service{
		httpClient:  &http.Client{},
		externalURL: externalAPIURL,
	}
}

// SendRiskNotification converts the RiskResult to JSON and sends it via HTTP POST.
func (s *Service) SendRiskNotification(ctx context.Context, riskData *pb.RiskResult) error {
	// 1. Convert the protobuf message to a map for easy JSON marshaling
	payload := map[string]interface{}{
		"trx_key":           riskData.GetTrxKey(),
		"risk_level":        riskData.GetRiskLevel().String(),
		"risk_score":        riskData.GetRiskScore(),
		"triggered_rules":   riskData.GetTriggeredRules(),
		"notification_type": "RISK_ALERT",
	}

	// 2. Marshal the map into a JSON byte slice
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	// 3. Create a new HTTP POST request
	req, err := http.NewRequestWithContext(ctx, "POST", s.externalURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 4. Send the request
	log.Printf("Sending notification to %s...", s.externalURL)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send http request: %w", err)
	}
	defer resp.Body.Close()

	// 5. Check the response
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("Successfully sent notification for TrxKey: %s, Status: %s", riskData.GetTrxKey(), resp.Status)
	} else {
		log.Printf("ERROR: Notification sent for TrxKey %s failed with status: %s", riskData.GetTrxKey(), resp.Status)
	}

	return nil
}
