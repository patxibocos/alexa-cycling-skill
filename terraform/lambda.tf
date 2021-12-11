data "archive_file" "alexa_cycling_lambda_code" {
  type        = "zip"
  source_file = "alexa-skill-lambda"
  output_path = "alexa-skill-lambda.zip"
}

data "archive_file" "appconfig_publisher_lambda_code" {
  type        = "zip"
  source_file = "appconfig-publisher-lambda"
  output_path = "appconfig-publisher-lambda.zip"
}

resource "aws_lambda_function" "alexa_cycling_lambda" {
  function_name    = "alexa_cycling_lambda"
  role             = aws_iam_role.alexa_cycling_role.arn
  runtime          = "go1.x"
  filename         = data.archive_file.alexa_cycling_lambda_code.output_path
  source_code_hash = data.archive_file.alexa_cycling_lambda_code.output_base64sha256
  handler          = "alexa-skill-lambda"
  publish          = true
  layers           = ["arn:aws:lambda:eu-west-3:493207061005:layer:AWS-AppConfig-Extension:46"]
  environment {
    variables = {
      AWS_APPCONFIG_URL = "http://localhost:2772/applications/alexa_cycling_appconfig/environments/alexa_cycling_appconfig_environment/configurations/alexa_cycling_appconfig_profile"
    }
  }
}

resource "aws_lambda_function" "appconfig_publisher_lambda" {
  function_name    = "appconfig_publisher_lambda"
  role             = aws_iam_role.appconfig_publisher_role.arn
  runtime          = "go1.x"
  filename         = data.archive_file.appconfig_publisher_lambda_code.output_path
  source_code_hash = data.archive_file.appconfig_publisher_lambda_code.output_base64sha256
  handler          = "appconfig-publisher-lambda"
  publish          = true
}

resource "aws_iam_role" "alexa_cycling_role" {
  name = "alexa_cycling_role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role" "appconfig_publisher_role" {
  name = "appconfig_publisher_role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_policy" "appconfig_get_configuration_policy" {
  name        = "appconfig_get_configuration_policy"
  description = "Read Appconfig policy"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "appconfig:GetConfiguration",
      "Effect": "Allow",
      "Resource": "arn:aws:appconfig:eu-west-3:654448679164:*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "test-attach" {
  role       = aws_iam_role.alexa_cycling_role.name
  policy_arn = aws_iam_policy.appconfig_get_configuration_policy.arn
}

resource "aws_iam_role_policy_attachment" "terraform_lambda_policy" {
  role       = aws_iam_role.alexa_cycling_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "another-policy" {
  role       = aws_iam_role.appconfig_publisher_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role" "alexa_cycling_appconfig_role" {
  name = "alexa_cycling_appconfig_role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "appconfig.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_policy" "start_deployment_policy" {
  name        = "start_deployment_policy"
  description = "Start deployment policy"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "appconfig:StartDeployment"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_policy" "read_s3_policy" {
  name        = "read_s3_policy"
  description = "Read S3 policy"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject", 
        "s3:GetObjectVersion"
      ],
      "Resource": [
        "arn:aws:s3:::alexacycling/cycling.data"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetBucketVersioning",
        "s3:GetBucketLocation",
        "s3:ListBucketVersions",
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::alexacycling"
      ]
    },
    {
      "Effect": "Allow",
      "Action": "s3:ListAllMyBuckets",
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "test-attach2" {
  role       = aws_iam_role.alexa_cycling_appconfig_role.name
  policy_arn = aws_iam_policy.read_s3_policy.arn
}

resource "aws_iam_role_policy_attachment" "another-attach" {
  role       = aws_iam_role.appconfig_publisher_role.name
  policy_arn = aws_iam_policy.read_s3_policy.arn
}

resource "aws_iam_role_policy_attachment" "another-attach2" {
  role       = aws_iam_role.appconfig_publisher_role.name
  policy_arn = aws_iam_policy.start_deployment_policy.arn
}

resource "aws_appconfig_application" "alexa_cycling_appconfig_application" {
  name        = "alexa_cycling_appconfig"
  description = "alexa_cycling_appconfig"
}

resource "aws_appconfig_configuration_profile" "alexa_cycling_appconfig_profile" {
  application_id     = aws_appconfig_application.alexa_cycling_appconfig_application.id
  description        = "alexa_cycling_appconfig_profile"
  name               = "alexa_cycling_appconfig_profile"
  location_uri       = "s3://alexacycling/cycling.data"
  retrieval_role_arn = aws_iam_role.alexa_cycling_appconfig_role.arn
}

resource "aws_appconfig_environment" "alexa_cycling_appconfig_environment" {
  name           = "alexa_cycling_appconfig_environment"
  description    = "alexa_cycling_appconfig_environment"
  application_id = aws_appconfig_application.alexa_cycling_appconfig_application.id
}

resource "aws_s3_bucket" "alexa_cycling_s3_bucket" {
  bucket = "alexacycling"
  acl    = "private"
  versioning {
    enabled = true
  }
}

resource "aws_lambda_permission" "allow_bucket" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.appconfig_publisher_lambda.arn
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.alexa_cycling_s3_bucket.arn
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = aws_s3_bucket.alexa_cycling_s3_bucket.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.appconfig_publisher_lambda.arn
    events              = ["s3:ObjectCreated:*"]
  }

  depends_on = [aws_lambda_permission.allow_bucket]
}

output "alexa_cycling_lambda_version" {
  value = aws_lambda_function.alexa_cycling_lambda.version
}

output "alexa_cycling_lambda_name" {
  value = aws_lambda_function.alexa_cycling_lambda.function_name
}

output "appconfig_publisher_lambda_version" {
  value = aws_lambda_function.appconfig_publisher_lambda.version
}

output "appconfig_publisher_lambda_name" {
  value = aws_lambda_function.appconfig_publisher_lambda.function_name
}