resource "aws_api_gateway_rest_api" "io_service" {
    name        = "io-service"
    body        = data.template_file.io_service.rendered
}

data "template_file" "io_service" {
    template = file("apiSpec.yaml")
}

resource "aws_api_gateway_deployment" "io_service" {
    depends_on = [aws_api_gateway_stage.io_service]
    rest_api_id = aws_api_gateway_rest_api.io_service.id
    stage_name = var.deployment_type
}