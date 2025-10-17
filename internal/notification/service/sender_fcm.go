package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	pb "github.com/wirsal/project-aegis/api/protos"
)

type fcmNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Image string `json:"image,omitempty"`
}

type fcmAndroidConfig struct {
	Priority string `json:"priority,omitempty"`
}

type fcmApnsConfig struct {
	Headers map[string]string `json:"headers,omitempty"`
}

type fcmRequestPayload struct {
	Tokens       []string         `json:"tokens"`
	Notification fcmNotification  `json:"notification"`
	Android      fcmAndroidConfig `json:"android,omitempty"`
	Apns         fcmApnsConfig    `json:"apns,omitempty"`
}

type FirebaseSender struct {
	httpClient *http.Client
	fcmURL     string
}

func NewFirebaseSender(url string) *FirebaseSender {
	return &FirebaseSender{
		httpClient: &http.Client{},
		fcmURL:     url,
	}
}

func (s *FirebaseSender) Send(ctx context.Context, req *pb.NotificationRequest) error {
	log.Printf("FirebaseSender: Preparing notification for recipient: %s", req.GetRecipient())

	payloadMap := req.GetPayload().AsMap()
	title, _ := payloadMap["title"].(string)
	body, _ := payloadMap["body"].(string)

	fcmPayload := fcmRequestPayload{
		Tokens: []string{req.GetRecipient()},
		Notification: fcmNotification{
			Title: title,
			Body:  body,
		},
		Android: fcmAndroidConfig{Priority: "normal"},
		Apns:    fcmApnsConfig{Headers: map[string]string{"apns-priority": "10"}},
	}

	jsonData, err := json.Marshal(fcmPayload)
	if err != nil {
		return fmt.Errorf("firebase sender: failed to marshal payload: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.fcmURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("firebase sender: failed to create request: %w", err)
	}

	// TODO: Tambahkan logic untuk mendapatkan dan menyisipkan OAuth2 Token di sini
	// httpReq.Header.Set("Authorization", "Bearer "+oauthToken)
	httpReq.Header.Set("Content-Type", "application/json")

	log.Printf("FirebaseSender: Sending POST request to %s...", s.fcmURL)
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		log.Printf("ERROR: FirebaseSender failed to send http request: %v", err)
		return fmt.Errorf("failed to send http request to FCM: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("FirebaseSender: Successfully sent notification, status: %s", resp.Status)
	} else {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Printf("ERROR: FirebaseSender received non-2xx status: %s, body: %s", resp.Status, string(bodyBytes))
		return fmt.Errorf("FCM returned non-2xx status: %s", resp.Status)
	}

	return nil
}
