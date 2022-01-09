#!/bin/bash
echo $FUNCTION_NAME
zip $FUNCTION_NAME alexa-skill-lambda
aws s3 cp $FUNCTION_NAME s3://$AWS_S3_BUCKET/