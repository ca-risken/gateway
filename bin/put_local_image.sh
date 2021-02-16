#!/bin/bash -e

cd "$(dirname "$0")"

# load env
. ../env.sh

# setting remote repository
TAG="local-test-$(date '+%Y%m%d')"
IMAGE_GATEWAY="gateway/gateway"
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query "Account" --output text)
REGISTORY="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"

# build & push
aws ecr get-login-password --region ${AWS_REGION} \
  | docker login \
    --username AWS \
    --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com

docker build --build-arg GITHUB_USER=${GITHUB_USER} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t ${IMAGE_GATEWAY}:${TAG} ..
docker tag ${IMAGE_GATEWAY}:${TAG} ${REGISTORY}/${IMAGE_GATEWAY}:${TAG}
docker push ${REGISTORY}/${IMAGE_GATEWAY}:${TAG}
