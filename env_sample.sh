#!/bin/bash -e

# github
export GITHUB_USER="your_name_here"
export GITHUB_TOKEN="your_token_here"

# GO
export GOPRIVATE="github.com/CyberAgent/*"

# grpc server
export FINDING_SVC_ADDR="finding:8001"
export IAM_SVC_ADDR="iam:8002"
export PROJECT_SVC_ADDR="project:8003"
export AWS_SVC_ADDR="aws:9001"
export OSINT_SVC_ADDR="osint:18080"