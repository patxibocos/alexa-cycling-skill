#!/bin/bash
ALIAS_NAME=$1

ALEXA_CYCLING_LAMBDA_VERSION=$(terraform output -raw alexa_cycling_lambda_version)
ALEXA_CYCLING_LAMBDA_NAME=$(terraform output -raw alexa_cycling_lambda_name)

aws lambda update-alias \
  --function-name $ALEXA_CYCLING_LAMBDA_NAME \
  --name $ALIAS_NAME \
  --function-version $ALEXA_CYCLING_LAMBDA_VERSION \
  > /dev/null \
|| \
aws lambda create-alias \
  --function-name $ALEXA_CYCLING_LAMBDA_NAME \
  --name $ALIAS_NAME \
  --function-version $ALEXA_CYCLING_LAMBDA_VERSION \
  > /dev/null