#!/bin/bash

echo "=== AWS EC2 Instance Manager Test Script ==="
echo

# Check if AWS credentials are configured
if ! aws sts get-caller-identity &>/dev/null; then
    echo "❌ AWS credentials not configured. Please run 'aws configure' first."
    exit 1
fi

echo "✅ AWS credentials configured"
echo

# Test dry-run list command
echo "🔍 Testing list command (dry-run)..."
./ec2-manager -action=list -dry-run
echo

# Test dry-run start command
echo "🚀 Testing start command (dry-run)..."
./ec2-manager -action=start -dry-run
echo

# Test dry-run stop command
echo "🛑 Testing stop command (dry-run)..."
./ec2-manager -action=stop -dry-run
echo

# Test with custom tags
echo "🏷️  Testing with custom tags (dry-run)..."
./ec2-manager -action=list -tag-key=Environment -tag-value=dev -dry-run
echo

echo "✅ All tests completed successfully!"
echo "💡 Remove -dry-run flag to perform actual operations"