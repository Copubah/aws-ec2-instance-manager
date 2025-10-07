.PHONY: build clean test run-list run-start run-stop build-lambda

# Default target
all: build

# Build the application
build:
	go build -o ec2-manager main.go

# Build for Lambda deployment
build-lambda:
	GOOS=linux GOARCH=amd64 go build -o ec2-manager main.go
	zip ec2-manager-lambda.zip ec2-manager

# Clean build artifacts
clean:
	rm -f ec2-manager ec2-manager-lambda.zip

# Run tests
test:
	go test -v ./...

# Quick run commands for testing
run-list:
	go run main.go -action=list -dry-run

run-start:
	go run main.go -action=start -dry-run

run-stop:
	go run main.go -action=stop -dry-run

# Install dependencies
deps:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run