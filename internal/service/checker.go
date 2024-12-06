package service

import (
	"context"
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"encoding/json"
	"fmt"
	"github.com/zinrai/redmine-ticket-sentinel/internal/domain"
	"math/rand"
)

type TicketChecker struct {
	issueRepo domain.IssueRepository
	userRepo  domain.UserRepository
	notifier  domain.Notifier
	cueValue  cue.Value
	ruleNames []string
}

type EvaluationResult struct {
	Message string
	Issue   *domain.Issue
	Mention string // String for Slack Mention
}

func NewTicketChecker(
	issueRepo domain.IssueRepository,
	userRepo domain.UserRepository,
	notifier domain.Notifier,
	cueDir string,
	ruleNames []string,
) (*TicketChecker, error) {
	// Load CUE configurations
	ctx := cuecontext.New()
	bis := load.Instances([]string{cueDir}, nil)
	if len(bis) == 0 {
		return nil, fmt.Errorf("no cue files found in %s", cueDir)
	}

	value := ctx.BuildInstance(bis[0])
	if value.Err() != nil {
		return nil, fmt.Errorf("failed to build cue instance: %w", value.Err())
	}

	return &TicketChecker{
		issueRepo: issueRepo,
		userRepo:  userRepo,
		notifier:  notifier,
		cueValue:  value,
		ruleNames: ruleNames,
	}, nil
}

func (c *TicketChecker) evaluateRule(ruleName string, issue *domain.Issue) (*EvaluationResult, error) {
	// Convert issue to JSON for CUE processing
	issueJSON, err := json.Marshal(issue)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal issue: %w", err)
	}

	// Create CUE context and process
	ctx := cuecontext.New()
	issueValue := ctx.CompileBytes(issueJSON)
	if issueValue.Err() != nil {
		return nil, fmt.Errorf("failed to compile issue: %w", issueValue.Err())
	}

	// Get rule from CUE
	rule := c.cueValue.LookupPath(cue.ParsePath(ruleName))
	if rule.Err() != nil {
		return nil, fmt.Errorf("rule not found: %w", rule.Err())
	}

	// Create evaluation context
	evalCtx := ctx.CompileString("_issue: _")
	evalCtx = evalCtx.FillPath(cue.ParsePath("_issue"), issueValue)

	// Evaluate rule
	result := rule.Unify(evalCtx)
	if result.Err() != nil {
		return nil, fmt.Errorf("evaluation failed: %w", result.Err())
	}

	// Extract evaluation result and message
	var evaluated bool
	var message, mention string

	if err := result.LookupPath(cue.ParsePath("evaluate")).Decode(&evaluated); err != nil {
		return nil, fmt.Errorf("failed to decode evaluation result: %w", err)
	}

	if !evaluated {
		return nil, nil
	}

	if err := result.LookupPath(cue.ParsePath("message")).Decode(&message); err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}

	if err := result.LookupPath(cue.ParsePath("mention")).Decode(&mention); err != nil {
		// If no mapping is found, use your Redmine login name
		mention = issue.Author.Login
	}

	return &EvaluationResult{
		Message: message,
		Issue:   issue,
		Mention: mention,
	}, nil
}

func (c *TicketChecker) CheckTickets(ctx context.Context) error {

	issues, err := c.issueRepo.GetIssues(ctx)
	if err != nil {
		return fmt.Errorf("failed to get issues: %w", err)
	}

	// Randomly select one rule
	ruleName := c.ruleNames[rand.Intn(len(c.ruleNames))]

	for _, issue := range issues {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		detailedIssue, err := c.issueRepo.GetIssueDetails(ctx, issue.ID)
		if err != nil {
			continue // Log error but continue with next issue
		}

		result, err := c.evaluateRule(ruleName, detailedIssue)
		if err != nil {
			continue // Log error but continue with next issue
		}

		if result != nil {
			if err := c.notifier.SendNotification(
				ctx,
				result.Mention,
				result.Issue,
				result.Message,
			); err != nil {
				continue // Log error but continue with next issue
			}
		}
	}

	return nil
}
