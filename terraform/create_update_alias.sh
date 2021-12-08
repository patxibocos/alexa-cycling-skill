#!/bin/bash
VERSION=$(cat output.json | jq -r '.lambda_version.value')
FUNCTION_NAME=$(cat output.json | jq -r '.lambda_name.value')

aws lambda update-alias \
  --function-name $FUNCTION_NAME \
  --name $GIT_BRANCH \
  --function-version $VERSION \
  > /dev/null \
|| \
aws lambda create-alias \
  --function-name $FUNCTION_NAME \
  --name $GIT_BRANCH \
  --function-version $VERSION \
  --description "The latest build in the $GIT_BRANCH branch"
  > /dev/null