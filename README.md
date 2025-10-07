# AWS EC2 Instance Manager

A Go-based automation tool for managing AWS EC2 instances using tags. This script can list, start, or stop EC2 instances based on configurable tag filters.

## Features

- ✅ List EC2 instances with specific tags
- ✅ Start/stop instances based on tags
- ✅ Configurable tag filtering
- ✅ Dry-run mode for safe testing
- ✅ Comprehensive error handling and logging
- ✅ Easy to extend for scheduling or Lambda deployment

## Prerequisites

- Go 1.19 or later
- AWS credentials configured (via AWS CLI, environment variables, or IAM roles)
- Appropriate IAM permissions for EC2 operations

## Required IAM Permissions

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeInstances",
                "ec2:StartInstances",
                "ec2:StopInstances"
            ],
            "Resource": "*"
        }
    ]
}
```

## Installation

1. Clone or download this repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

### Basic Commands

```bash
# List all instances with AutoManage=true tag
go run main.go -action=list

# Start instances with AutoManage=true tag
go run main.go -action=start

# Stop instances with AutoManage=true tag
go run main.go -action=stop

# Use custom tag filter
go run main.go -action=list -tag-key=Environment -tag-value=dev

# Dry run (safe testing)
go run main.go -action=stop -dry-run

# Specify different region
go run main.go -action=list -region=us-west-2
```

### Command Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `-action` | `list` | Action to perform: `list`, `start`, `stop` |
| `-tag-key` | `AutoManage` | Tag key to filter instances |
| `-tag-value` | `true` | Tag value to filter instances |
| `-region` | `us-east-1` | AWS region |
| `-dry-run` | `false` | Perform a dry run without making changes |

## Building

```bash
# Build for current platform
go build -o ec2-manager main.go

# Build for Linux (Lambda deployment)
GOOS=linux GOARCH=amd64 go build -o ec2-manager main.go
```

## Example Output

```
[EC2-Manager] 2024/01/15 10:30:00 Starting EC2 instance management - Action: list, Tag: AutoManage=true, Region: us-east-1
[EC2-Manager] 2024/01/15 10:30:01 Found 3 instances with tag AutoManage=true
[EC2-Manager] 2024/01/15 10:30:01 Listing instances:
Instance ID          State           Type            Name                          
--------------------------------------------------------------------------------
i-1234567890abcdef0  running         t3.micro        web-server-1                  
i-0987654321fedcba0  stopped         t3.small        worker-node-1                 
i-abcdef1234567890   running         t3.medium       database-server               
```

## Extension Ideas

### 1. Scheduling Support
Add cron-like scheduling functionality:
- Parse schedule configuration
- Implement time-based triggers
- Add timezone support

### 2. Lambda Deployment
Package for AWS Lambda:
- Create Lambda handler wrapper
- Add CloudWatch Events integration
- Environment-based configuration

### 3. Additional Features
- Instance health checks before operations
- Batch operations with rate limiting
- Slack/SNS notifications
- Cost optimization reports
- Multi-region support

## Configuration File

Copy `config.example.json` to `config.json` and customize for your environment. This enables future extensions like scheduling and Lambda deployment.

## Error Handling

The script includes comprehensive error handling for:
- AWS API errors
- Network connectivity issues
- Invalid configurations
- Missing permissions
- Instance state conflicts

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details