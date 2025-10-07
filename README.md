# AWS EC2 Instance Manager
- A production-ready Go application for automated AWS EC2 instance management using tag-based filtering. Built with AWS SDK v2, this tool provides safe, reliable operations for starting, stopping, and listing EC2 instances.

## Features

- List EC2 instances with specific tags
- Start/stop instances based on configurable tag filters
- Dry-run mode for safe testing and validation
- Comprehensive error handling and structured logging
- Multi-region support
- Lambda deployment ready
- Production-grade security and best practices

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   CLI Client    │    │  Lambda Function │    │  CloudWatch     │
│                 │    │                  │    │  Events         │
└─────────┬───────┘    └─────────┬────────┘    └─────────┬───────┘
          │                      │                       │
          │              ┌───────▼───────┐               │
          └──────────────►│  EC2 Manager  │◄──────────────┘
                         │               │
                         └───────┬───────┘
                                 │
                    ┌────────────▼────────────┐
                    │      AWS EC2 API        │
                    │                         │
                    │  ┌─────────────────┐    │
                    │  │   Instance 1    │    │
                    │  │ Tag: AutoManage │    │
                    │  └─────────────────┘    │
                    │                         │
                    │  ┌─────────────────┐    │
                    │  │   Instance 2    │    │
                    │  │ Tag: AutoManage │    │
                    │  └─────────────────┘    │
                    │                         │
                    │  ┌─────────────────┐    │
                    │  │   Instance N    │    │
                    │  │ Tag: AutoManage │    │
                    │  └─────────────────┘    │
                    └─────────────────────────┘

Flow:
1. Client/Lambda triggers EC2 Manager with action (list/start/stop)
2. EC2 Manager queries AWS EC2 API with tag filters
3. Manager performs requested action on filtered instances
4. Results logged and returned to caller
```

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

## Project Structure

```
aws-ec2-instance-manager/
├── main.go                 # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   │   └── config.go
│   └── ec2/               # EC2 operations
│       ├── manager.go
│       └── manager_test.go
├── .github/
│   └── workflows/
│       └── ci.yml         # GitHub Actions CI/CD
├── config.example.json    # Configuration template
├── lambda_handler.go      # Lambda deployment wrapper
├── deploy-lambda.sh       # Lambda deployment script
├── test-script.sh         # Testing utilities
├── Makefile              # Build automation
├── SECURITY.md           # Security guidelines
└── LICENSE               # MIT License
```

## Environment Variables

The application supports configuration via environment variables:

```bash
export AWS_REGION=us-east-1
export EC2_TAG_KEY=AutoManage
export EC2_TAG_VALUE=true
```

## Example Output

```
[EC2-Manager] 2024/01/15 10:30:00 main.go:45: Starting EC2 instance management - Action: list, Tag: AutoManage=true, Region: us-east-1
[EC2-Manager] 2024/01/15 10:30:01 manager.go:67: Found 3 instances with tag AutoManage=true
[EC2-Manager] 2024/01/15 10:30:01 manager.go:89: Listing instances:
Instance ID          State           Type            Name                          
--------------------------------------------------------------------------------
i-1234567890abcdef0  running         t3.micro        web-server-1                  
i-0987654321fedcba0  stopped         t3.small        worker-node-1                 
i-abcdef1234567890   running         t3.medium       database-server               
```

## Best Practices Implemented

### Code Organization
- Modular package structure with clear separation of concerns
- Comprehensive error handling with wrapped errors
- Structured logging with file and line information
- Input validation and configuration management

### Security
- Environment variable support for sensitive configuration
- IAM permission validation and least privilege examples
- Secure credential handling (no hardcoded secrets)
- Security policy documentation

### Testing & CI/CD
- Unit tests with table-driven test patterns
- GitHub Actions workflow for automated testing
- Multi-platform build support
- Security scanning with Gosec
- Code coverage reporting

### Production Readiness
- Graceful error handling and recovery
- Comprehensive logging for debugging
- Dry-run mode for safe operations
- Resource cleanup and proper context handling

## Extension Ideas

### 1. Scheduling Support
```go
// Add to config.json
{
  "schedule": {
    "enabled": true,
    "start_time": "08:00",
    "stop_time": "18:00",
    "timezone": "UTC"
  }
}
```

### 2. Multi-Region Operations
```bash
# Support multiple regions
./ec2-manager -action=list -regions=us-east-1,us-west-2,eu-west-1
```

### 3. Advanced Filtering
```bash
# Multiple tag filters
./ec2-manager -action=list -tags="Environment=prod,Team=backend"
```

## Deployment Options

### 1. Binary Deployment
```bash
# Download from releases
curl -L https://github.com/username/aws-ec2-instance-manager/releases/latest/download/ec2-manager-linux-amd64 -o ec2-manager
chmod +x ec2-manager
```

### 2. Lambda Deployment
```bash
# Use provided deployment script
./deploy-lambda.sh
```

### 3. Container Deployment
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o ec2-manager main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/ec2-manager .
CMD ["./ec2-manager"]
```

## Contributing

We welcome contributions! Please see our contributing guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes with tests
4. Run the test suite (`make test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development Setup
```bash
# Clone the repository
git clone https://github.com/username/aws-ec2-instance-manager.git
cd aws-ec2-instance-manager

# Install dependencies
go mod download

# Run tests
make test

# Build locally
make build
```

## Security

Please review our [Security Policy](SECURITY.md) for information about:
- Reporting vulnerabilities
- Security best practices
- Secure configuration examples

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- Create an issue for bug reports or feature requests
- Check existing issues before creating new ones
- Provide detailed information for faster resolution