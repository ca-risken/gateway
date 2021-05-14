#!/bin/bash -e

# github
export GITHUB_USER="your_name_here"
export GITHUB_TOKEN="your_token_here"

# GO
export GOPRIVATE="github.com/CyberAgent/*"

# build
export AWS_REGION="ap-northeast-1"

# gateway
export PORT="8000"
export DEBUG="true"
export USER_IDENTITY_HEADER="x-amzn-oidc-identity"
export OIDC_DATA_HEADER="x-amzn-oidc-data"
export IDP_PROVIDER_NAME="YOUR_IDP1,YOUR_IDP2"

# grpc server
export FINDING_SVC_ADDR="finding:8001"
export IAM_SVC_ADDR="iam:8002"
export PROJECT_SVC_ADDR="project:8003"
export ALERT_SVC_ADDR="alert:8004"
export REPORT_SVC_ADDR="report:8005"
export AWS_SVC_ADDR="aws:9001"
export AWS_ACTIVITY_SVC_ADDR="activity:9007"
export OSINT_SVC_ADDR="osint:18081"
export DIAGNOSIS_SVC_ADDR="diagnosis:19001"
export CODE_SVC_ADDR="code:10001"
export GOOGLE_SVC_ADDR="google:11001"
