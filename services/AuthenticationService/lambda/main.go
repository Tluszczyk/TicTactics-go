package main

import (
	"services/lib"
	"services/lib/types"

	"services/AuthenticationService/cmd"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := cmd.HandleRequest(types.Request{
		HTTPMethod: request.HTTPMethod,
		Body:       request.Body,
	})

	return lib.CreateApiResponse(
		response.StatusCode,
		response.Body,
	), err
}

func main() {
	lambda.Start(handler)
}
