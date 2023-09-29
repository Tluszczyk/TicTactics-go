resource "aws_api_gateway_rest_api" "io_service" {
    name        = "io-service"
    body        = data.template_file.io_service.rendered
}

data "template_file" "io_service" {
    template = file("apiSpec/apiSpec.yaml")

    vars = {
        user_management_lambda_uri = aws_lambda_function.user_management_service.invoke_arn
        authentication_lambda_uri = aws_lambda_function.authentication_service.invoke_arn
        game_management_lambda_uri = aws_lambda_function.game_management_service.invoke_arn
    }
}

resource "aws_api_gateway_deployment" "io_service" {
    depends_on = [aws_api_gateway_rest_api.io_service]
    rest_api_id = aws_api_gateway_rest_api.io_service.id
    stage_name = var.deployment_type
}

output "api_gateway_id" {
    value = aws_api_gateway_rest_api.io_service.id
}

output "api_gateway_endpoint" {
    value = format(
        "http://%s.execute-api.localhost.localstack.cloud:4566/%s",
        aws_api_gateway_rest_api.io_service.id,
        aws_api_gateway_deployment.io_service.stage_name
    )   
}