#!/bin/bash

# AWS EC2 Manager Lambda Deployment Script

set -e

FUNCTION_NAME="ec2-auto-manager"
REGION="us-east-1"
ROLE_NAME="ec2-manager-lambda-role"

echo "Deploying EC2 Manager to AWS Lambda..."

# Build for Lambda
echo "Building Lambda package..."
GOOS=linux GOARCH=amd64 go build -o ec2-manager main.go
zip ec2-manager-lambda.zip ec2-manager

# Check if IAM role exists
if ! aws iam get-role --role-name $ROLE_NAME &>/dev/null; then
    echo "Creating IAM role..."
    
    # Create trust policy
    cat > trust-policy.json << EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "lambda.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
EOF

    # Create IAM role
    aws iam create-role \
        --role-name $ROLE_NAME \
        --assume-role-policy-document file://trust-policy.json

    # Attach basic Lambda execution policy
    aws iam attach-role-policy \
        --role-name $ROLE_NAME \
        --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole

    # Create and attach EC2 policy
    cat > ec2-policy.json << EOF
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
EOF

    aws iam put-role-policy \
        --role-name $ROLE_NAME \
        --policy-name EC2ManagementPolicy \
        --policy-document file://ec2-policy.json

    echo "Waiting for IAM role to propagate..."
    sleep 10
fi

# Get role ARN
ROLE_ARN=$(aws iam get-role --role-name $ROLE_NAME --query 'Role.Arn' --output text)

# Deploy or update Lambda function
if aws lambda get-function --function-name $FUNCTION_NAME &>/dev/null; then
    echo "Updating existing Lambda function..."
    aws lambda update-function-code \
        --function-name $FUNCTION_NAME \
        --zip-file fileb://ec2-manager-lambda.zip
else
    echo "Creating new Lambda function..."
    aws lambda create-function \
        --function-name $FUNCTION_NAME \
        --runtime go1.x \
        --role $ROLE_ARN \
        --handler ec2-manager \
        --zip-file fileb://ec2-manager-lambda.zip \
        --timeout 300 \
        --memory-size 128 \
        --description "Automated EC2 instance management"
fi

# Clean up temporary files
rm -f trust-policy.json ec2-policy.json

echo "PASS: Lambda function deployed successfully!"
echo "Function name: $FUNCTION_NAME"
echo "Region: $REGION"
echo "Role ARN: $ROLE_ARN"
echo
echo "To test the function:"
echo "aws lambda invoke --function-name $FUNCTION_NAME --payload '{\"action\":\"list\",\"dry_run\":true}' response.json"