package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/zinrai/redmine-ticket-sentinel/internal/infrastructure"
	"github.com/zinrai/redmine-ticket-sentinel/internal/service"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Redmine struct {
		BaseURL   string `yaml:"base_url"`
		Username  string `yaml:"username"`
		Password  string `yaml:"password"`
		ProjectID string `yaml:"project_id"`
		QueryID   int    `yaml:"query_id"`
	} `yaml:"redmine"`

	Slack struct {
		WebhookURL string `yaml:"webhook_url"`
	} `yaml:"slack"`

	Rules struct {
		Path  string   `yaml:"path"`
		Names []string `yaml:"names"`
	} `yaml:"rules"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	debugFlag := flag.Bool("debug", false, "Use stdout for notifications instead of Slack")
	ruleFlag := flag.String("rule", "", "Specify a single rule to evaluate (e.g., #StoryWithoutChildren)")
	flag.Parse()

	// Setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load configuration
	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Redmine client
	redmineClient := infrastructure.NewRedmineClient(
		config.Redmine.BaseURL,
		config.Redmine.Username,
		config.Redmine.Password,
		config.Redmine.ProjectID,
		config.Redmine.QueryID,
	)

	// Choose notifier based on debug flag
	var notifier domain.Notifier
	if *debugFlag {
		notifier = infrastructure.NewStdoutNotifier()
	} else {
		notifier = infrastructure.NewSlackNotifier(
			config.Slack.WebhookURL,
			config.Redmine.BaseURL,
		)
	}

	// Use specified rule or all rules from config
	ruleNames := config.Rules.Names
	if *ruleFlag != "" {
		log.Printf("Using specified rule: %s", *ruleFlag)
		ruleNames = []string{*ruleFlag}
	}

	// Initialize ticket checker
	checker, err := service.NewTicketChecker(
		redmineClient,
		notifier,
		config.Rules.Path,
		ruleNames,
	)
	if err != nil {
		log.Fatalf("Failed to initialize ticket checker: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Run ticket checks
	if err := checker.CheckTickets(ctx); err != nil {
		log.Printf("Error checking tickets: %v", err)
		os.Exit(1)
	}
}
