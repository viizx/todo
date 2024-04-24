package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type SMSError struct {
	To      string
	Message string
	Err     error
}

type InfobipSMSService struct {
	APIKey string
	Sender string
	URL    string
	Client *http.Client
	Errors chan SMSError
}

func (s *InfobipSMSService) HandleErrors() {
	for err := range s.Errors {
		log.Printf("Failed to send SMS to %s: %v", err.To, err.Err)
	}
}

func NewInfobipSMSService(apiKey, sender, urlStr string) *InfobipSMSService {
	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		log.Fatalf("Invalid URL provided: %v", err)
	}
	return &InfobipSMSService{
		APIKey: apiKey,
		Sender: sender,
		URL:    parsedURL.String(),
		Client: &http.Client{},
		Errors: make(chan SMSError, 100),
	}
}

func (s *InfobipSMSService) SendSMS(ctx context.Context, to, message string) error {
	payload := map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"from": s.Sender,
				"destinations": []map[string]string{
					{"to": to},
				},
				"text": message,
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		s.Errors <- SMSError{To: to, Message: message, Err: err}
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.URL, bytes.NewBuffer(body))
	if err != nil {
		s.Errors <- SMSError{To: to, Message: message, Err: err}
		return err
	}

	req.Header.Add("Authorization", "App "+s.APIKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := s.Client.Do(req)
	fmt.Println(res.Body)
	if err != nil {
		s.Errors <- SMSError{To: to, Message: message, Err: err}
		return err
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		s.Errors <- SMSError{To: to, Message: message, Err: err}
		return err
	}

	if res.StatusCode >= 300 {
		s.Errors <- SMSError{To: to, Message: string(responseBody), Err: fmt.Errorf("received non-success status code %d", res.StatusCode)}
		return fmt.Errorf("HTTP error: %d %s", res.StatusCode, responseBody)
	}

	fmt.Println(string(responseBody))
	return nil
}
