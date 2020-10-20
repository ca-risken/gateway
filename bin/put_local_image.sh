#!/bin/bash -e

cd "$(dirname "$0")"

# load env
. ../env.sh

# setting remote repository
export IMAGE="gateway/gateway"
export TAG="local-test-$(date '+%Y%m%d')"

# build & push
docker build --build-arg GITHUB_USER=${GITHUB_USER} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t ${IMAGE}:${TAG} ..
$(aws ecr get-login --no-include-email --region ${AWS_DEFAULT_REGION})
docker tag ${IMAGE}:${TAG} ${REGISTORY}/${IMAGE}:${TAG}
docker push ${REGISTORY}/${IMAGE}:${TAG}
