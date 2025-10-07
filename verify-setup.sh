#!/bin/bash

echo "AWS EC2 Instance Manager - Setup Verification"
echo "============================================="

# Check Go installation
echo "Checking Go installation..."
if command -v go &> /dev/null; then
    echo "PASS: Go is installed: $(go version)"
else
    echo "FAIL: Go is not installed"
    exit 1
fi

# Check if we can build the project
echo "Building project..."
if go build -o ec2-manager main.go; then
    echo "PASS: Build successful"
else
    echo "FAIL: Build failed"
    exit 1
fi

# Run tests
echo "Running tests..."
if go test ./...; then
    echo "PASS: All tests passed"
else
    echo "FAIL: Tests failed"
    exit 1
fi

# Check if binary works
echo "Testing binary..."
if ./ec2-manager -h > /dev/null 2>&1; then
    echo "PASS: Binary works correctly"
else
    echo "FAIL: Binary execution failed"
    exit 1
fi

# Verify file structure
echo "Verifying project structure..."
required_files=(
    "main.go"
    "internal/config/config.go"
    "internal/ec2/manager.go"
    "internal/ec2/manager_test.go"
    "README.md"
    "LICENSE"
    "SECURITY.md"
    ".gitignore"
    "Makefile"
)

for file in "${required_files[@]}"; do
    if [[ -f "$file" ]]; then
        echo "PASS: $file exists"
    else
        echo "FAIL: $file missing"
        exit 1
    fi
done

# Check git status
echo "Git status..."
if git status --porcelain | grep -q .; then
    echo "WARNING: Uncommitted changes detected"
    git status --short
else
    echo "PASS: All changes committed"
fi

echo ""
echo "Setup verification complete!"
echo ""
echo "Next steps:"
echo "1. Create GitHub repository (see github-setup.md)"
echo "2. Push code: git push -u origin main"
echo "3. Configure AWS credentials for testing"
echo "4. Test with: ./test-script.sh"
echo ""
echo "Repository features:"
echo "- Production-ready Go application"
echo "- Modular architecture with best practices"
echo "- Comprehensive test coverage"
echo "- Security policies and documentation"
echo "- Lambda deployment support"
echo "- Multi-platform build support"