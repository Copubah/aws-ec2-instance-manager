// +build lambda

package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// LambdaEvent represents the input event for Lambda
type LambdaEvent struct {
	Action   string `json:"action"`
	TagKey   string `json:"tag_key,omitempty"`
	TagValue string `json:"tag_value,omitempty"`
	Region   string `json:"region,omitempty"`
	DryRun   bool   `json:"dry_run,omitempty"`
}

// LambdaResponse represents the Lambda response
type LambdaResponse struct {
	StatusCode int               `json:"statusCode"`
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
}

// HandleLambdaRequest processes Lambda events
func HandleLambdaRequest(ctx context.Context, event events.CloudWatchEvent) (LambdaResponse, error) {
	// Parse the event detail
	var lambdaEvent LambdaEvent
	if err := json.Unmarshal(event.Detail, &lambdaEvent); err != nil {
		return LambdaResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Error parsing event: %v", err),
		}, nil
	}

	// Set defaults
	if lambdaEvent.TagKey == "" {
		lambdaEvent.TagKey = "AutoManage"
	}
	if lambdaEvent.TagValue == "" {
		lambdaEvent.TagValue = "true"
	}
	if lambdaEvent.Region == "" {
		lambdaEvent.Region = "us-east-1"
	}
	if lambdaEvent.Action == "" {
		lambdaEvent.Action = "list"
	}

	// Create config
	cfg := &Config{
		Action:   lambdaEvent.Action,
		TagKey:   lambdaEvent.TagKey,
		TagValue: lambdaEvent.TagValue,
		Region:   lambdaEvent.Region,
		DryRun:   lambdaEvent.DryRun,
	}

	// Initialize EC2 manager
	manager, err := NewEC2Manager(cfg.Region)
	if err != nil {
		return LambdaResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Failed to initialize EC2 manager: %v", err),
		}, nil
	}

	// Execute the action
	if err := manager.ManageInstances(cfg); err != nil {
		return LambdaResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Failed to manage instances: %v", err),
		}, nil
	}

	return LambdaResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("Successfully executed action: %s", cfg.Action),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

// Lambda entry point (uncomment when building for Lambda)
// func main() {
// 	lambda.Start(HandleLambdaRequest)
// }