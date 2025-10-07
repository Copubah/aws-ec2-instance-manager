# Security Policy

## Supported Versions

We actively support the following versions with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability, please follow these steps:

1. **Do not** create a public GitHub issue for security vulnerabilities
2. Email the maintainers directly with details of the vulnerability
3. Include steps to reproduce the issue
4. Allow reasonable time for the issue to be addressed before public disclosure

## Security Best Practices

When using this tool:

### AWS Credentials
- Never hardcode AWS credentials in your code
- Use IAM roles when running on EC2 instances
- Use AWS CLI profiles or environment variables for local development
- Regularly rotate access keys

### IAM Permissions
- Follow the principle of least privilege
- Only grant necessary EC2 permissions
- Use resource-based policies when possible
- Regularly audit IAM permissions

### Network Security
- Use VPC endpoints for AWS API calls when possible
- Ensure proper security group configurations
- Monitor CloudTrail logs for API activity

### Application Security
- Always use dry-run mode for testing
- Validate input parameters
- Use structured logging (avoid logging sensitive data)
- Keep dependencies updated

## Secure Configuration

### Environment Variables
```bash
# Recommended environment setup
export AWS_REGION=us-east-1
export EC2_TAG_KEY=AutoManage
export EC2_TAG_VALUE=true
```

### IAM Policy Example
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeInstances"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "ec2:StartInstances",
                "ec2:StopInstances"
            ],
            "Resource": "arn:aws:ec2:*:*:instance/*",
            "Condition": {
                "StringEquals": {
                    "ec2:ResourceTag/AutoManage": "true"
                }
            }
        }
    ]
}
```

This policy ensures that the application can only manage instances with the `AutoManage=true` tag.