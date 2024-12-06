package domain

import "context"

type IssueRepository interface {
	GetIssues(ctx context.Context) ([]Issue, error)
	GetIssueDetails(ctx context.Context, issueID int) (*Issue, error)
}

type UserRepository interface {
	GetUsers(ctx context.Context) (map[int]string, error)
}

type Notifier interface {
	SendNotification(ctx context.Context, username string, issue *Issue, message string) error
}
