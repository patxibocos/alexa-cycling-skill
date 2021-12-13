#!/bin/bash
GIT_HASH=$(git rev-parse --short "$GITHUB_SHA")

ALEXA_CYCLING_LAMBDA_VERSION=$(terraform output -raw alexa_cycling_lambda_version)
ALEXA_CYCLING_LAMBDA_NAME=$(terraform output -raw alexa_cycling_lambda_name)

aws lambda update-alias \
  --function-name $ALEXA_CYCLING_LAMBDA_NAME \
  --name $GIT_HASH \
  --function-version $ALEXA_CYCLING_LAMBDA_VERSION \
  > /dev/null \
|| \
aws lambda create-alias \
  --function-name $ALEXA_CYCLING_LAMBDA_NAME \
  --name $GIT_HASH \
  --function-version $ALEXA_CYCLING_LAMBDA_VERSION \
  --description "Build from $GIT_HASH commit"
  > /dev/null