package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds the application configuration
type Config struct {
	Action    string
	TagKey    string
	TagValue  string
	Region    string
	DryRun    bool
}

// Validate ensures the configuration is valid
func (c *Config) Validate() error {
	validActions := []string{"list", "start", "stop"}
	for _, valid := range validActions {
		if c.Action == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid action: %s. Valid actions: %s", c.Action, strings.Join(validActions, ", "))
}

// LoadFromEnv loads configuration from environment variables
func (c *Config) LoadFromEnv() {
	if region := os.Getenv("AWS_REGION"); region != "" {
		c.Region = region
	}
	if tagKey := os.Getenv("EC2_TAG_KEY"); tagKey != "" {
		c.TagKey = tagKey
	}
	if tagValue := os.Getenv("EC2_TAG_VALUE"); tagValue != "" {
		c.TagValue = tagValue
	}
}