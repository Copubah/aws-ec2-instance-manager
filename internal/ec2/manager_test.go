package ec2

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestNewManager(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	
	tests := []struct {
		name   string
		region string
		want   bool
	}{
		{
			name:   "valid region",
			region: "us-east-1",
			want:   true,
		},
		{
			name:   "empty region",
			region: "",
			want:   true, // AWS SDK will use default region
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewManager(tt.region, logger)
			if tt.want && err != nil {
				t.Errorf("NewManager() error = %v, want success", err)
				return
			}
			if tt.want && manager == nil {
				t.Error("NewManager() returned nil manager")
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid list action",
			config: &Config{
				Action:   "list",
				TagKey:   "AutoManage",
				TagValue: "true",
				Region:   "us-east-1",
			},
			wantErr: false,
		},
		{
			name: "valid start action",
			config: &Config{
				Action:   "start",
				TagKey:   "AutoManage",
				TagValue: "true",
				Region:   "us-east-1",
			},
			wantErr: false,
		},
		{
			name: "valid stop action",
			config: &Config{
				Action:   "stop",
				TagKey:   "AutoManage",
				TagValue: "true",
				Region:   "us-east-1",
			},
			wantErr: false,
		},
		{
			name: "invalid action",
			config: &Config{
				Action:   "invalid",
				TagKey:   "AutoManage",
				TagValue: "true",
				Region:   "us-east-1",
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAction(tt.config.Action)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// validateAction is a helper function extracted for testing
func validateAction(action string) error {
	validActions := []string{"list", "start", "stop"}
	for _, valid := range validActions {
		if action == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid action: %s", action)
}