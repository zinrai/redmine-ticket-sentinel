package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zinrai/redmine-ticket-sentinel/internal/domain"
	"net/http"
	"time"
)

type SlackNotifier struct {
	webhookURL string
	baseURL    string
	httpClient *http.Client
}

func NewSlackNotifier(webhookURL, baseURL string) *SlackNotifier {
	return &SlackNotifier{
		webhookURL: webhookURL,
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (n *SlackNotifier) SendNotification(ctx context.Context, mention string, issue *domain.Issue, message string) error {
	title := fmt.Sprintf("%s (%s): %s", issue.Project.Name, issue.Status.Name, issue.Subject)
	content := fmt.Sprintf("*%s*\n<@%s> %s/issues/%d\n%s",
		title, mention, n.baseURL, issue.ID, message)

	payload := struct {
		Text string `json:"text"`
	}{
		Text: content,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	time.Sleep(time.Second) // Rate limiting
	return nil
}
