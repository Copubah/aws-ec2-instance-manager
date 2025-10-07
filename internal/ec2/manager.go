package ec2

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Manager handles EC2 operations
type Manager struct {
	client *ec2.Client
	logger *log.Logger
}

// Config represents the configuration for EC2 operations
type Config struct {
	Action    string
	TagKey    string
	TagValue  string
	Region    string
	DryRun    bool
}

// NewManager creates a new EC2 manager instance
func NewManager(region string, logger *log.Logger) (*Manager, error) {
	ctx := context.Background()
	
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}
	
	client := ec2.NewFromConfig(cfg)
	
	return &Manager{
		client: client,
		logger: logger,
	}, nil
}

// ManageInstances performs the requested action on filtered instances
func (m *Manager) ManageInstances(ctx context.Context, cfg *Config) error {
	m.logger.Printf("Starting EC2 instance management - Action: %s, Tag: %s=%s, Region: %s", 
		cfg.Action, cfg.TagKey, cfg.TagValue, cfg.Region)
	
	instances, err := m.getInstancesByTag(ctx, cfg.TagKey, cfg.TagValue)
	if err != nil {
		return fmt.Errorf("failed to get instances: %w", err)
	}
	
	if len(instances) == 0 {
		m.logger.Printf("No instances found with tag %s=%s", cfg.TagKey, cfg.TagValue)
		return nil
	}
	
	m.logger.Printf("Found %d instances with tag %s=%s", len(instances), cfg.TagKey, cfg.TagValue)
	
	switch cfg.Action {
	case "list":
		return m.listInstances(instances)
	case "start":
		return m.startInstances(ctx, instances, cfg.DryRun)
	case "stop":
		return m.stopInstances(ctx, instances, cfg.DryRun)
	default:
		return fmt.Errorf("unsupported action: %s", cfg.Action)
	}
}

func (m *Manager) getInstancesByTag(ctx context.Context, tagKey, tagValue string) ([]types.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:" + tagKey),
				Values: []string{tagValue},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running", "stopped", "stopping", "pending"},
			},
		},
	}
	
	result, err := m.client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %w", err)
	}
	
	var instances []types.Instance
	for _, reservation := range result.Reservations {
		instances = append(instances, reservation.Instances...)
	}
	
	return instances, nil
}

func (m *Manager) listInstances(instances []types.Instance) error {
	m.logger.Println("Listing instances:")
	fmt.Printf("%-20s %-15s %-15s %-30s\n", "Instance ID", "State", "Type", "Name")
	fmt.Println(strings.Repeat("-", 80))
	
	for _, instance := range instances {
		name := m.getInstanceName(instance)
		fmt.Printf("%-20s %-15s %-15s %-30s\n", 
			aws.ToString(instance.InstanceId),
			string(instance.State.Name),
			string(instance.InstanceType),
			name)
	}
	
	return nil
}

func (m *Manager) startInstances(ctx context.Context, instances []types.Instance, dryRun bool) error {
	var stoppedInstances []string
	
	for _, instance := range instances {
		if instance.State.Name == types.InstanceStateNameStopped {
			stoppedInstances = append(stoppedInstances, aws.ToString(instance.InstanceId))
		}
	}
	
	if len(stoppedInstances) == 0 {
		m.logger.Println("No stopped instances to start")
		return nil
	}
	
	if dryRun {
		m.logger.Printf("DRY RUN: Would start %d instances: %s", 
			len(stoppedInstances), strings.Join(stoppedInstances, ", "))
		return nil
	}
	
	input := &ec2.StartInstancesInput{
		InstanceIds: stoppedInstances,
	}
	
	result, err := m.client.StartInstances(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to start instances: %w", err)
	}
	
	m.logger.Printf("Successfully initiated start for %d instances", len(result.StartingInstances))
	for _, instance := range result.StartingInstances {
		m.logger.Printf("Instance %s: %s -> %s", 
			aws.ToString(instance.InstanceId),
			string(instance.PreviousState.Name),
			string(instance.CurrentState.Name))
	}
	
	return nil
}

func (m *Manager) stopInstances(ctx context.Context, instances []types.Instance, dryRun bool) error {
	var runningInstances []string
	
	for _, instance := range instances {
		if instance.State.Name == types.InstanceStateNameRunning {
			runningInstances = append(runningInstances, aws.ToString(instance.InstanceId))
		}
	}
	
	if len(runningInstances) == 0 {
		m.logger.Println("No running instances to stop")
		return nil
	}
	
	if dryRun {
		m.logger.Printf("DRY RUN: Would stop %d instances: %s", 
			len(runningInstances), strings.Join(runningInstances, ", "))
		return nil
	}
	
	input := &ec2.StopInstancesInput{
		InstanceIds: runningInstances,
	}
	
	result, err := m.client.StopInstances(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to stop instances: %w", err)
	}
	
	m.logger.Printf("Successfully initiated stop for %d instances", len(result.StoppingInstances))
	for _, instance := range result.StoppingInstances {
		m.logger.Printf("Instance %s: %s -> %s", 
			aws.ToString(instance.InstanceId),
			string(instance.PreviousState.Name),
			string(instance.CurrentState.Name))
	}
	
	return nil
}

func (m *Manager) getInstanceName(instance types.Instance) string {
	for _, tag := range instance.Tags {
		if aws.ToString(tag.Key) == "Name" {
			return aws.ToString(tag.Value)
		}
	}
	return "N/A"
}