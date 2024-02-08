# User Management Service

resource "aws_lambda_function" "user_management_service" {
  filename          = format("%s/UserManagementLambda.zip", var.files_path)
  function_name     = "UserManagementLambda"
  role              = aws_iam_role.user_management_service.arn
  handler           = "bootstrap"
  timeout           = 30
  memory_size       = 512
  runtime           = "provided.al2"

  environment {
    variables = {
      USERS_TABLE_NAME = aws_dynamodb_table.this.name,
      PASSWORD_HASHES_TABLE_NAME = aws_dynamodb_table.this.name,
      ENDPOINT_URL = var.endpoint_url
    }
  }
}

resource "aws_iam_role" "user_management_service" {
  name = "UserManagementLambdaRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = "sts:AssumeRole",
      Effect = "Allow",
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}


# Authentication Service

resource "aws_lambda_function" "authentication_service" {
  filename          = format("%s/AuthenticationLambda.zip", var.files_path)
  function_name     = "AuthenticationLambda"
  role              = aws_iam_role.authentication_service.arn
  handler           = "bootstrap"
  timeout           = 30
  memory_size       = 512
  runtime           = "provided.al2"
}

resource "aws_iam_role" "authentication_service" {
  name = "AuthenticationLambdaRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = "sts:AssumeRole",
      Effect = "Allow",
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}


# Game Management Service

resource "aws_lambda_function" "game_management_service" {
  filename          = format("%s/GameManagementLambda.zip", var.files_path)
  function_name     = "GameManagementLambda"
  role              = aws_iam_role.game_management_service.arn
  handler           = "bootstrap"
  timeout           = 30
  memory_size       = 512
  runtime           = "provided.al2"
}

resource "aws_iam_role" "game_management_service" {
  name = "GameManagementLambdaRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = "sts:AssumeRole",
      Effect = "Allow",
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}