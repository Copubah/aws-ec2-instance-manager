package main

import (
	"context"
	"flag"
	"log"
	"os"

	"go-aws-ec2-automation/internal/config"
	"go-aws-ec2-automation/internal/ec2"
)

func main() {
	cfg := parseFlags()
	
	// Load environment variables
	cfg.LoadFromEnv()
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}
	
	// Initialize logger
	logger := log.New(os.Stdout, "[EC2-Manager] ", log.LstdFlags|log.Lshortfile)
	
	// Create EC2 manager
	manager, err := ec2.NewManager(cfg.Region, logger)
	if err != nil {
		log.Fatalf("Failed to initialize EC2 manager: %v", err)
	}

	// Create context with timeout
	ctx := context.Background()
	
	// Convert to internal config format
	ec2Config := &ec2.Config{
		Action:   cfg.Action,
		TagKey:   cfg.TagKey,
		TagValue: cfg.TagValue,
		Region:   cfg.Region,
		DryRun:   cfg.DryRun,
	}

	if err := manager.ManageInstances(ctx, ec2Config); err != nil {
		log.Fatalf("Failed to manage instances: %v", err)
	}
}

func parseFlags() *config.Config {
	cfg := &config.Config{}
	
	flag.StringVar(&cfg.Action, "action", "list", "Action to perform: list, start, stop")
	flag.StringVar(&cfg.TagKey, "tag-key", "AutoManage", "Tag key to filter instances")
	flag.StringVar(&cfg.TagValue, "tag-value", "true", "Tag value to filter instances")
	flag.StringVar(&cfg.Region, "region", "us-east-1", "AWS region")
	flag.BoolVar(&cfg.DryRun, "dry-run", false, "Perform a dry run without making changes")
	
	flag.Parse()
	
	return cfg
}

