package infrastructure

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zinrai/redmine-ticket-sentinel/internal/domain"
)

type RedmineClient struct {
	baseURL    string
	username   string
	password   string
	projectID  string
	queryID    int
	httpClient *http.Client
}

func NewRedmineClient(baseURL, username, password, projectID string, queryID int) *RedmineClient {
	return &RedmineClient{
		baseURL:    baseURL,
		username:   username,
		password:   password,
		projectID:  projectID,
		queryID:    queryID,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *RedmineClient) createAuthHeader() string {
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.password)))
	return fmt.Sprintf("Basic %s", auth)
}

func (c *RedmineClient) makeRequest(ctx context.Context, url string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.createAuthHeader())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func (c *RedmineClient) GetIssues(ctx context.Context) ([]domain.Issue, error) {
	var result struct {
		Issues []domain.Issue `json:"issues"`
	}

	offset := 0
	var allIssues []domain.Issue

	for {
		url := fmt.Sprintf("%s/issues.json?query_id=%d&project_id=%s&limit=100&offset=%d",
			c.baseURL, c.queryID, c.projectID, offset)

		if err := c.makeRequest(ctx, url, &result); err != nil {
			return nil, err
		}

		allIssues = append(allIssues, result.Issues...)

		if len(result.Issues) < 100 {
			break
		}
		offset += 100
	}

	return allIssues, nil
}

func (c *RedmineClient) GetIssueDetails(ctx context.Context, issueID int) (*domain.Issue, error) {
	var result struct {
		Issue domain.Issue `json:"issue"`
	}

	url := fmt.Sprintf("%s/issues/%d.json?include=children", c.baseURL, issueID)
	if err := c.makeRequest(ctx, url, &result); err != nil {
		return nil, err
	}

	return &result.Issue, nil
}

func (c *RedmineClient) GetUsers(ctx context.Context) (map[int]string, error) {
	var result struct {
		Users []domain.User `json:"users"`
	}

	users := make(map[int]string)
	offset := 0

	for {
		url := fmt.Sprintf("%s/users.json?limit=100&offset=%d", c.baseURL, offset)
		if err := c.makeRequest(ctx, url, &result); err != nil {
			return nil, err
		}

		for _, user := range result.Users {
			users[user.ID] = user.Login
		}

		if len(result.Users) < 100 {
			break
		}
		offset += 100
	}

	return users, nil
}
