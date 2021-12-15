variable "ALEXA_SKILL_ID" {}
variable "AWS_S3_BUCKET" {}
variable "AWS_S3_OBJECT_KEY" {}

data "archive_file" "alexa_cycling_lambda_code" {
  type        = "zip"
  source_file = "alexa-skill-lambda"
  output_path = "alexa-skill-lambda.zip"
}

resource "aws_lambda_function" "alexa_cycling_lambda" {
  function_name    = "alexa_cycling_lambda"
  role             = aws_iam_role.alexa_cycling_role.arn
  runtime          = "go1.x"
  filename         = data.archive_file.alexa_cycling_lambda_code.output_path
  source_code_hash = data.archive_file.alexa_cycling_lambda_code.output_base64sha256
  handler          = "alexa-skill-lambda"
  publish          = true
  environment {
    variables = {
      AWS_S3_BUCKET     = var.AWS_S3_BUCKET
      AWS_S3_OBJECT_KEY = var.AWS_S3_OBJECT_KEY
      ALEXA_SKILL_ID    = var.ALEXA_SKILL_ID
    }
  }
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
  inline_policy {
    name = "s3_get_object_policy"
    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [{
        Effect   = "Allow",
        Action   = ["s3:GetObject"]
        Resource = "arn:aws:s3:::${var.AWS_S3_BUCKET}/${var.AWS_S3_OBJECT_KEY}"
      }]
    })
  }
}

resource "aws_iam_role_policy_attachment" "terraform_lambda_policy" {
  role       = aws_iam_role.alexa_cycling_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_s3_bucket" "alexa_cycling_s3_bucket" {
  bucket = var.AWS_S3_BUCKET
  acl    = "private"
}

resource "aws_lambda_permission" "allow_alexa_skill" {
  statement_id       = "AllowExecutionFromAlexa"
  action             = "lambda:InvokeFunction"
  function_name      = aws_lambda_function.alexa_cycling_lambda.arn
  principal          = "alexa-appkit.amazon.com"
  event_source_token = var.ALEXA_SKILL_ID
}

output "alexa_cycling_lambda_version" {
  value = aws_lambda_function.alexa_cycling_lambda.version
}

output "alexa_cycling_lambda_name" {
  value = aws_lambda_function.alexa_cycling_lambda.function_name
}