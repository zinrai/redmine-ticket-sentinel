package infrastructure

import (
	"context"
	"fmt"
	"github.com/zinrai/redmine-ticket-sentinel/internal/domain"
)

type StdoutNotifier struct{}

func NewStdoutNotifier() *StdoutNotifier {
	return &StdoutNotifier{}
}

func (n *StdoutNotifier) SendNotification(ctx context.Context, mention string, issue *domain.Issue, message string) error {
	fmt.Printf("\n=== Ticket Information ===\n")
	fmt.Printf("Rule Message: %s\n", message)
	fmt.Printf("Issue ID: %d\n", issue.ID)
	fmt.Printf("Subject: %s\n", issue.Subject)
	fmt.Printf("Status: %s (ID: %d)\n", issue.Status.Name, issue.Status.ID)
	fmt.Printf("Project: %s\n", issue.Project.Name)
	fmt.Printf("Tracker: %s (ID: %d)\n", issue.Tracker.Name, issue.Tracker.ID)
	fmt.Printf("Assignee: %s (ID: %d)\n", issue.AssignedTo.Login, issue.AssignedTo.ID)
	fmt.Printf("Author: %s (ID: %d)\n", issue.Author.Login, issue.Author.ID)
	fmt.Printf("Children Count: %d\n", len(issue.Children))
	if len(issue.Children) > 0 {
		fmt.Printf("\n--- Children ---\n")
		for _, child := range issue.Children {
			fmt.Printf("Child ID: %d, Subject: %s, Status: %s, Tracker: %s\n",
				child.ID, child.Subject, child.Status.Name, child.Tracker.Name)
		}
	}
	fmt.Printf("\nSlack Mention: @%s\n", mention)
	fmt.Printf("=======================\n\n")
	return nil
}
